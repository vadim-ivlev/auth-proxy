var app = require('express')();
var bodyParser = require('body-parser');

app.use(bodyParser.json()); // for parsing application/json
app.use(bodyParser.urlencoded({ extended: true })); // for parsing application/x-www-form-urlencoded

app.post('/*', handler)
app.get('/*', handler)
app.options('/*', handler)

app.listen(3000);


// F U N C S -----------------------------------------
function handler(req, res) {
    console.log('----------------------------')
    console.log('method:',req.method)
    console.log('headers:',req.headers)
    console.log('url:',req.url)
    console.log('originalUrl:',req.originalUrl)
    console.log('path:',req.path)
    console.log('query:',req.query)
    console.log('body:', req.body)
    console.log('')
    res.send("response")
    // res.end()
}

