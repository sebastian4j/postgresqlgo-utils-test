package postgresqlgoutilstest

import (
	"bufio"
	"context"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/sebastian4j/postgresqlgo"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func Setup(t *testing.T) testcontainers.Container {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:10-alpine",
		ExposedPorts: []string{"5432"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
		Env: map[string]string{
			"POSTGRES_USER":             "postgres",
			"POSTGRES_PASSWORD":         "postgres",
			"POSTGRES_DB":               "postgres",
			"POSTGRES_HOST_AUTH_METHOD": "trust",
		},
	}
	rc, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Error(err)
	}
	mappedPort, _ := rc.MappedPort(ctx, "5432")
	hostIP, _ := rc.Host(ctx)
	log.Printf("%s:%s", hostIP, mappedPort.Port())
	os.Setenv("POSTGRES_HOST", hostIP)
	os.Setenv("POSTGRES_PORT", mappedPort.Port())
	os.Setenv("POSTGRES_USER", "postgres")
	os.Setenv("POSTGRES_PASSWORD", "postgres")
	os.Setenv("POSTGRES_DB", "postgres")
	return rc
}

func LoadScript(file *os.File, p *postgresqlgo.Postgresqlgo) bool {
	scanner := bufio.NewScanner(file)
	var sb strings.Builder
	for scanner.Scan() {
		sb.WriteString(scanner.Text())
		sb.WriteString("\n")
	}
	sqls := strings.Split(sb.String(), ";")
	conn := p.Conn()
	for _, cmd := range sqls {
		x := strings.Trim(cmd, "")
		if x != "" {
			_, err := conn.Exec(context.Background(), x)
			if err != nil {
				log.Fatalf("error al enviar comando a postgres %v %s", err, cmd)
			}

		}
	}
	conn.Release()
	return true
}

// LoadScriptChan carga el contenido de un archivo en posgres agrupando las lineas a enviar hasta completar q
func LoadScriptChan(file *os.File, p *postgresqlgo.Postgresqlgo, q int) int {
	if q < 2 {
		log.Fatalf("'q' tiene que ser sobre 1 y fue: %d", q)
	}
	cuenta := 0
	scanner := bufio.NewScanner(file)
	var sb strings.Builder
	for scanner.Scan() {
		sb.WriteString(scanner.Text())
		sb.WriteString("\n")
	}
	sqls := strings.Split(sb.String(), ";")

	var dones []chan bool
	acc := []string{}
	for _, cmd := range sqls {
		x := strings.Trim(cmd, "")
		if x != "" {
			acc = append(acc, x)
			if len(acc) >= q {
				done := make(chan bool)
				dones = append(dones, done)
				go func(xx []string) {
					conn, err := p.ConnWithErr()
					for err != nil {
						conn, err = p.ConnWithErr()
					}
					cuenta++
					batch := &pgx.Batch{}
					for _, q := range xx {
						batch.Queue(q)
					}
					br := conn.SendBatch(context.Background(), batch)
					ct, err := br.Exec()
					if err != nil {
						log.Fatalf("error al enviar batch %v", err)
					}
					if ct.RowsAffected() == int64(0) {
						log.Fatalf("error en insert batch: %d", ct.RowsAffected())
					}
					conn.Release()
					done <- true
				}(acc)
				acc = []string{}
			}
		}
	}
	for _, b := range dones {
		<-b
	}
	return len(dones)
}
