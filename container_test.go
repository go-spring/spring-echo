/*
 * Copyright 2012-2019 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package SpringEcho_test

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/go-spring/spring-core/assert"
	"github.com/go-spring/spring-core/web"
	"github.com/go-spring/spring-echo"
	"github.com/labstack/echo"
)

func TestContext_PanicEchoHttpError(t *testing.T) {
	c := SpringEcho.NewContainer(web.ContainerConfig{Port: 8080})
	c.GetMapping("/", func(ctx web.Context) {
		panic(echo.ErrTooManyRequests)
	})
	go c.Start()
	defer c.Stop(context.Background())
	time.Sleep(10 * time.Millisecond)
	response, err := http.Get("http://127.0.0.1:8080/")
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	b, _ := ioutil.ReadAll(response.Body)
	fmt.Println(response.Status, string(b))
	assert.Equal(t, response.StatusCode, http.StatusTooManyRequests)
	assert.Equal(t, string(b), `Too Many Requests`)
}

func TestContext_PanicString(t *testing.T) {
	c := SpringEcho.NewContainer(web.ContainerConfig{Port: 8080})
	c.GetMapping("/", func(ctx web.Context) {
		panic("this is an error")
	})
	go c.Start()
	defer c.Stop(context.Background())
	time.Sleep(10 * time.Millisecond)
	response, err := http.Get("http://127.0.0.1:8080/")
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	b, _ := ioutil.ReadAll(response.Body)
	fmt.Println(response.Status, string(b))
	assert.Equal(t, response.StatusCode, http.StatusOK)
	assert.Equal(t, string(b), "\"this is an error\"")
}

func TestContext_PanicError(t *testing.T) {
	c := SpringEcho.NewContainer(web.ContainerConfig{Port: 8080})
	c.GetMapping("/", func(ctx web.Context) {
		panic(errors.New("this is an error"))
	})
	go c.Start()
	defer c.Stop(context.Background())
	time.Sleep(10 * time.Millisecond)
	response, err := http.Get("http://127.0.0.1:8080/")
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	b, _ := ioutil.ReadAll(response.Body)
	fmt.Println(response.Status, string(b))
	assert.Equal(t, response.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, string(b), `this is an error`)
}

func TestContext_PanicWebHttpError(t *testing.T) {
	c := SpringEcho.NewContainer(web.ContainerConfig{Port: 8080})
	c.GetMapping("/", func(ctx web.Context) {
		panic(&web.HttpError{
			Code:    http.StatusNotFound,
			Message: http.StatusText(http.StatusNotFound),
		})
	})
	go c.Start()
	defer c.Stop(context.Background())
	time.Sleep(10 * time.Millisecond)
	response, err := http.Get("http://127.0.0.1:8080/")
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	b, _ := ioutil.ReadAll(response.Body)
	fmt.Println(response.Status, string(b))
	assert.Equal(t, response.StatusCode, http.StatusNotFound)
	assert.Equal(t, string(b), `Not Found`)
}

type dummyFilter struct{}

func (f *dummyFilter) Invoke(ctx web.Context, chain web.FilterChain) {
	panic(&web.HttpError{
		Code:    http.StatusMethodNotAllowed,
		Message: http.StatusText(http.StatusMethodNotAllowed),
	})
}

func TestFilter_PanicWebHttpError(t *testing.T) {
	c := SpringEcho.NewContainer(web.ContainerConfig{Port: 8080})
	c.GetMapping("/", func(ctx web.Context) {
		ctx.String("OK!")
	}, &dummyFilter{})
	go c.Start()
	defer c.Stop(context.Background())
	time.Sleep(10 * time.Millisecond)
	response, err := http.Get("http://127.0.0.1:8080/")
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	b, _ := ioutil.ReadAll(response.Body)
	fmt.Println(response.Status, string(b))
	assert.Equal(t, response.StatusCode, http.StatusMethodNotAllowed)
	assert.Equal(t, string(b), `Method Not Allowed`)
}
