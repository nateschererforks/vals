package ssm

import (
	"fmt"
	"os"
	"sort"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/google/go-cmp/cmp"

	"github.com/nateschererforks/vals/pkg/config"
	"github.com/nateschererforks/vals/pkg/log"
)

type mockedSSM struct {
	ssmiface.SSMAPI
	Error  awserr.Error
	Output *ssm.GetParametersByPathOutput
	Path   string
}

func Output(params map[string]string) *ssm.GetParametersByPathOutput {
	ssmParams := []*ssm.Parameter{}

	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		ssmParams = append(ssmParams, &ssm.Parameter{
			Name:  aws.String(k),
			Value: aws.String(params[k]),
		})
	}

	return &ssm.GetParametersByPathOutput{
		Parameters: ssmParams,
	}
}

func (m mockedSSM) GetParametersByPath(in *ssm.GetParametersByPathInput) (*ssm.GetParametersByPathOutput, error) {
	path := *in.Path
	if path != m.Path {
		return nil, fmt.Errorf("unexpected path: %s", path)
	}

	return m.Output, m.Error
}

func (m mockedSSM) GetParametersByPathPages(in *ssm.GetParametersByPathInput, fn func(o *ssm.GetParametersByPathOutput, lastPage bool) bool) error {
	path := *in.Path
	if path != m.Path {
		return fmt.Errorf("unexpected path: %s", path)
	}

	fn(m.Output, true)
	return m.Error
}

func TestGetStringMap(t *testing.T) {
	cases := []struct {
		ssm       mockedSSM
		want      map[string]interface{}
		key       string
		wantErr   string
		recursive bool
	}{
		// {
		// 	key:     "bar",
		// 	wantErr: "ssm: get parameters by path: ParameterNotFound: parameter not found\ncaused by: simulated parameter-not-found error",
		// 	ssm: mockedSSM{
		// 		Path:   "/bar",
		// 		Output: nil,
		// 		Error:  awserr.New(ssm.ErrCodeParameterNotFound, "parameter not found", errors.New("simulated parameter-not-found error")),
		// 	},
		// },
		// {
		// 	key: "foo",
		// 	want: map[string]interface{}{
		// 		"bar": "BAR",
		// 		"baz": "BAZ",
		// 	},
		// 	ssm: mockedSSM{
		// 		Path: "/foo",
		// 		Output: Output(map[string]string{
		// 			"/foo/bar": `BAR`,
		// 			"/foo/baz": `BAZ`,
		// 		}),
		// 		Error: nil,
		// 	},
		// },
		{
			key:       "foo",
			recursive: true,
			want: map[string]interface{}{
				"bar": map[string]interface{}{
					"a": "A",
					"b": "B",
				},
				"baz": "BAZ",
			},
			ssm: mockedSSM{
				Path: "/foo",
				Output: Output(map[string]string{
					"/foo/bar":   `BAR`,
					"/foo/bar/a": `A`,
					"/foo/bar/b": `B`,
					"/foo/baz":   `BAZ`,
				}),
				Error: nil,
			},
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			conf := map[string]interface{}{}

			if c.recursive {
				conf["recursive"] = "true"
			}

			p := New(log.New(log.Config{Output: os.Stderr}), config.MapConfig{M: conf})

			p.ssmClient = c.ssm

			got, err := p.GetStringMap(c.key)

			if err != nil {
				if err.Error() != c.wantErr {
					t.Fatalf("unexpected error: want %q, got %q", c.wantErr, err.Error())
				}
			} else {
				if c.wantErr != "" {
					t.Fatalf("expected error did not occur: want %q, got none", c.wantErr)
				}
			}

			if diff := cmp.Diff(c.want, got); diff != "" {
				t.Errorf("unexpected result: -(want), +(got)\n%s", diff)
			}
		})
	}
}
