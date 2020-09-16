package storage

import (
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
)

func logJSON(in interface{}) {
	data, _ := jsoniter.MarshalIndent(in, "", "    ")
	storage.logger.Printf("%s\n", string(data))
}
func TestStorage(t *testing.T) {
	var storage = New(&Config{
		RootDir: "temp/test2",
	})

	f, _ := storage.Create("test22.csv")
	f.WriteString("asdfasdfasdfasdfsadf")

	f1, _ := storage.Open("test22.csv")
	f1.WriteString("\n12312312312312")

	storage.Create("test23.csv")

	storage.Create("images/1.png")
	storage.Create("images/2.png")
	storage.Create("images/3.png")

	file, err := storage.Stat("test22.csv")
	assert.NoError(t, err)
	storage.logger.Printf("%v\n", file)

	files, err := storage.Walk()
	assert.NoError(t, err)
	logJSON(files)

	storage.Remove("test23.csv")

	storage.RemoveAll("images")
}
