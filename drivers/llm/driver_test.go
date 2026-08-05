package llm_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	llmdrivers "github.com/usesnipet/snipet/drivers/llm"
)

func TestRegistryRegistersProviders(t *testing.T) {
	reg := llmdrivers.Registry()
	want := []string{"openai", "groq", "ollama", "mistral", "openrouter"}
	require.ElementsMatch(t, want, reg.Names())

	for _, id := range want {
		driver, ok := reg.Get(id)
		require.True(t, ok, id)
		info := driver.Info()
		require.NotEmpty(t, info.Name)
		require.NotNil(t, info.ConfigurationSchema)
		props, ok := info.ConfigurationSchema["properties"].(map[string]any)
		require.True(t, ok, id)
		require.Contains(t, props, "model")
		if id == "ollama" {
			require.Contains(t, props, "endpoint")
		} else {
			require.Contains(t, props, "api_key")
		}
	}

	_, ok := reg.Get("gemini")
	require.False(t, ok)
}
