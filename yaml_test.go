package deepmerge_test

import (
	"errors"
	"testing"

	"github.com/TwiN/deepmerge"
	"gopkg.in/yaml.v3"
)

func TestYAML(t *testing.T) {
	scenarios := []struct {
		name        string
		config      deepmerge.Config
		dst         string
		src         string
		expected    string
		expectedErr error
	}{
		{
			name:        "invalid-dst",
			dst:         `wat`,
			src:         ``,
			expected:    ``,
			expectedErr: errors.New("yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `wat` into map[string]interface {}"),
		},
		{
			name:        "invalid-src",
			dst:         ``,
			src:         `wat`,
			expected:    ``,
			expectedErr: errors.New("yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `wat` into map[string]interface {}"),
		},
		{
			name: "simple-endpoint-merge",
			dst: `endpoints: 
  - name: one
    url: https://example.com
    client:
      timeout: 5s
    conditions:
      - "[CONNECTED] == true"
      - "[STATUS] == 200"
    alerts:
      - type: slack
        failure-threshold: 5

  - name: two
    url: https://example.org
    conditions:
      - "len([BODY]) > 0"`,
			src: `endpoints: 
  - name: three
    url: https://twin.sh/health
    conditions:
      - "[STATUS] == 200"
      - "[BODY].status == UP"`,
			expected: `endpoints: 
  - name: one
    url: https://example.com
    client:
      timeout: 5s
    conditions:
      - "[CONNECTED] == true"
      - "[STATUS] == 200"
    alerts:
      - type: slack
        failure-threshold: 5

  - name: two
    url: https://example.org
    conditions:
      - "len([BODY]) > 0"

  - name: three
    url: https://twin.sh/health
    conditions:
      - "[STATUS] == 200"
      - "[BODY].status == UP"
`,
		},
		{
			name: "deep-merge-with-map-slice-and-primitive",
			dst: `
metrics: true 

alerting:
  slack:
    webhook-url: https://hooks.slack.com/services/xxx/yyy/zzz
    default-alert:
      description: "health check failed"
      send-on-resolved: true
      failure-threshold: 5
      success-threshold: 5

endpoints:
  - name: example
    url: https://example.org
    interval: 5s`,
			src: `
debug: true

alerting:
  discord:
    webhook-url: https://discord.com/api/webhooks/xxx/yyy

endpoints:
  - name: frontend
    url: https://example.com`,
			expected: `
metrics: true
debug: true

alerting:
  discord:
    webhook-url: https://discord.com/api/webhooks/xxx/yyy
  slack:
    webhook-url: https://hooks.slack.com/services/xxx/yyy/zzz
    default-alert:
      description: "health check failed"
      send-on-resolved: true
      failure-threshold: 5
      success-threshold: 5

endpoints:
  - interval: 5s
    name: example
    url: https://example.org

  - name: frontend
    url: https://example.com
`,
		},
		{ // only maps and slices can be merged. If there are duplicate keys that have a primitive value, then that's an error.
			name:        "duplicate-key-with-primitive-value",
			config:      deepmerge.Config{PreventMultipleDefinitionsOfKeysWithPrimitiveValue: true}, // NOTE: true is the default
			dst:         `metrics: true`,
			src:         `metrics: false`,
			expectedErr: deepmerge.ErrKeyWithPrimitiveValueDefinedMoreThanOnce,
		},
		{
			name:   "duplicate-key-with-primitive-value-with-preventDuplicateKeysWithPrimitiveValue-set-to-false",
			config: deepmerge.Config{PreventMultipleDefinitionsOfKeysWithPrimitiveValue: false},
			dst: `metrics: true
debug: true`,
			src: `metrics: false`,
			expected: `metrics: false
debug: true`,
		},
		{
			name: "readme-example",
			dst: `debug: true
client:
  insecure: true
users:
  - id: 1
    firstName: John
    lastName: Doe
  - id: 2
    firstName: Jane
    lastName: Doe`,
			src: `client:
  timeout: 5s
users:
  - id: 3
    firstName: Bob
    lastName: Smith`,
			expected: `
client:
  insecure: true
  timeout: 5s
debug: true
users:
  - firstName: John
    id: 1
    lastName: Doe
  - firstName: Jane
    id: 2
    lastName: Doe
  - firstName: Bob
    id: 3
    lastName: Smith`,
		},
	}
	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			output, err := deepmerge.YAML([]byte(scenario.dst), []byte(scenario.src), scenario.config)
			if !errors.Is(err, scenario.expectedErr) && !(scenario.expectedErr != nil && err.Error() == scenario.expectedErr.Error()) {
				t.Errorf("[%s] expected error %v, got %v", scenario.name, scenario.expectedErr, err)
			}
			// Just so we don't have to worry about the formatting, we'll unmarshal the output and marshal it again.
			expectedAsMap, outputAsMap := make(map[string]interface{}), make(map[string]interface{})
			if len(output) > 0 {
				if err := yaml.Unmarshal(output, &outputAsMap); err != nil {
					t.Errorf("[%s] failed to unmarshal output: %v", scenario.name, err)
				}
			}
			if len(scenario.expected) > 0 {
				if err := yaml.Unmarshal([]byte(scenario.expected), &expectedAsMap); err != nil {
					t.Errorf("[%s] failed to unmarshal expected: %v", scenario.name, err)
				}
			}
			formattedOutput, err := yaml.Marshal(outputAsMap)
			if err != nil {
				t.Errorf("[%s] should've been able to re-marshal output: %v", scenario.name, err)
			}
			formattedExpected, err := yaml.Marshal(expectedAsMap)
			if err != nil {
				t.Errorf("[%s] should've been able to re-marshal expected: %v", scenario.name, err)
			}
			// Compare what we got vs what we expected
			if string(formattedOutput) != string(formattedExpected) {
				t.Errorf("[%s] expected:\n%s\n\ngot:\n%s", scenario.name, string(formattedExpected), string(formattedOutput))
			}
		})
	}
}
