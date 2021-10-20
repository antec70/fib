package internal

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"go-fib/app/config"
	"log"
	"net/http"
	"strconv"
	"time"
)

var ctx = context.Background()

type Server struct {
	config     config.ParamsLocal
	rdb        *redis.Client
	serverHttp *http.Server
}

func NewServer(c config.ParamsLocal) *Server {
	return &Server{
		config: c,
		rdb: redis.NewClient(&redis.Options{
			Addr:     c.Redis.Adr,
			Password: c.Redis.Password,
			DB:       c.Redis.Db,
		})}
}

func (s *Server) RPing() {
	_, er := s.rdb.Ping(ctx).Result()
	if er != nil {
		log.Fatal(er)
	}
}

func (s *Server) Start() error {
	s.routerConf()
	s.RPing()
	httpServer := &http.Server{
		Addr:    ":" + s.config.Port,
		Handler: s.routerConf(),
	}
	s.serverHttp = httpServer
	fmt.Println("Start api")
	return httpServer.ListenAndServe()
}
func (s *Server) routerConf() http.Handler {
	r := mux.NewRouter()
	r.Use(commonMiddleware)
	r.HandleFunc("/fib", s.handleIndex()).Methods("GET")
	return r

}
func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func (s *Server) handleIndex() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var x uint32
		var y uint32
		keysX, ok := r.URL.Query()["x"]
		if !ok || len(keysX[0]) < 1 {
			ResponseBad("missing x value", w)
			return
		}

		keysY, ok := r.URL.Query()["y"]
		if !ok || len(keysY[0]) < 1 {
			ResponseBad("missing y value", w)
			return
		}

		u64X, er := strconv.ParseUint(keysX[0], 10, 32)
		if er != nil {
			ResponseBad("X - value error: "+er.Error(), w)
			return
		}
		x = uint32(u64X)

		u64Y, er := strconv.ParseUint(keysY[0], 10, 32)
		if er != nil {
			ResponseBad("Y - value error: "+er.Error(), w)
			return
		}

		y = uint32(u64Y)
		f := s.CalcFib(x, y)
		ResponseOk(f, w)
	}
}
func Fib(n uint32) uint32 {
	if n <= 1 {
		return n
	}
	var n2, n1 uint32 = 0, 1

	for i := uint32(2); i < n; i++ {
		n2, n1 = n1, n1+n2
	}
	return n2 + n1
}

func SlowFib(n uint32) uint32 {
	if n == 0 {
		return 0
	}
	if n < 2 {
		return 1
	}
	return SlowFib(n-2) + SlowFib(n-1)
}
func (s *Server) CalcFib(x, y uint32) []FibItem {

	var res []FibItem
	for i := x; i < y; i++ {
		r := FibItem{
			Position: i,
			Item:     Fib(i),
		}
		v, er := s.rdb.Get(ctx, fmt.Sprint(i)).Result()
		if er == redis.Nil {
			r.Item = Fib(i)
			s.rdb.Set(ctx, fmt.Sprint(i), r.Item, 0)
		} else if er != nil {
			fmt.Errorf("%s", er)
			r.Item = Fib(i)
		} else {
			in, _ := strconv.Atoi(v)
			r.Item = uint32(in)
		}

		res = append(res, r)
	}
	return res
}

func (s *Server) Shutdown() error {
	fmt.Println("Shutdown api...")
	ctx2, _ := context.WithTimeout(context.Background(), time.Duration(2)*time.Second)
	err := s.serverHttp.Shutdown(ctx2)
	return err
}
