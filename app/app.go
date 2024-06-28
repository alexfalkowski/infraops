package app

import (
	"bytes"
	"errors"
	"io"
	"net/http"

	"github.com/avast/retry-go/v4"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// ErrVersionMismatch for app.
var ErrVersionMismatch = errors.New("version mismatch")

type (
	// App to be created.
	App struct {
		Memory        Memory
		ID            string
		Name          string
		Domain        string
		Version       string
		ConfigVersion string
		InitVersion   string
		SecretVolumes []string
	}

	// Memory for apps.
	Memory struct {
		Min string
		Max string
	}

	createFn func(ctx *pulumi.Context, app *App) error
)

// CreateApp in the cluster.
func CreateApp(ctx *pulumi.Context, app *App) error {
	fns := []createFn{
		createServiceAccount, createNetworkPolicy,
		createConfigMap, createPodDisruptionBudget,
		createDeployment, createService, createIngress,
	}

	for _, fn := range fns {
		if err := fn(ctx, app); err != nil {
			return err
		}
	}

	return nil
}

// Probe the app.
func (a *App) Probe(ctx *pulumi.Context, path, body string) (string, error) {
	if ctx.DryRun() {
		return "", nil
	}

	fn := func() ([]byte, error) {
		req, err := http.NewRequestWithContext(ctx.Context(), "POST", a.url(path), bytes.NewReader([]byte(body)))
		if err != nil {
			return nil, err
		}

		req.Header.Set("Content-Type", "application/json")

		res, err := a.client().Do(req)
		if err != nil {
			return nil, err
		}

		defer res.Body.Close()

		if res.Header.Get("Service-Version") != a.Version {
			return nil, ErrVersionMismatch
		}

		return io.ReadAll(res.Body)
	}
	fnc := func(err error) bool {
		return errors.Is(err, ErrVersionMismatch)
	}

	d, err := retry.DoWithData(fn, retry.RetryIf(fnc))

	return string(d), err
}

func (a *App) url(path string) string {
	return "https://" + a.Name + "." + a.Domain + "/" + path
}

func (a *App) client() *http.Client {
	return http.DefaultClient
}
