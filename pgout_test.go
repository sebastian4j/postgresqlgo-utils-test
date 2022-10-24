package postgresqlgoutilstest

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/sebastian4j/postgresqlgo"
	"github.com/stretchr/testify/assert"
)

func TestPostgresConnection(t *testing.T) {
	p := postgresqlgo.Postgresqlgo{}
	x := Setup(t)
	defer x.Terminate(context.Background())
	t.Run("puedo generar conexiones", func(t *testing.T) {
		for a := 0; a < 1000; a++ {
			con := p.Conn()
			con.Release()
			con, err := p.ConnWithErr()
			if err != nil {
				t.Fatal(err)
			}
			con.Release()
		}
	})

	t.Run("puedo cargar un archivo secuencial", func(t *testing.T) {
		f, err := os.OpenFile("./testdata/bd.sql", os.O_RDONLY, 0644)
		if err != nil {
			log.Fatal("error al abrir archivo", err)
		}
		assert.Equal(t, true, LoadScript(f, &p))
		f.Close()
	})

	t.Run("puedo cargar un archivo en paralelo", func(t *testing.T) {
		f, err := os.OpenFile("./testdata/chan.sql", os.O_RDONLY, 0644)
		if err != nil {
			log.Fatal("error al abrir archivo", err)
		}
		assert.Equal(t, 5, LoadScriptChan(f, &p, 2))
		f.Close()
	})

}
