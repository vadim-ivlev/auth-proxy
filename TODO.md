## TODO:

- When proxying to url with path http://abc.com/a/b/c. Remove prefix?
- Do routing by hands `apps.Any("/:appname/*proxypath", Proxy) -> apps.Any("*proxypath", Proxy)`
