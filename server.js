var app = require('express')();
var bodyParser = require('body-parser');

app.use(bodyParser.json()); // for parsing application/json
app.use(bodyParser.urlencoded({ extended: true })); // for parsing application/x-www-form-urlencoded

// app.post('/*', handler)
// app.get('/*', handler)
// app.options('/*', handler)
app.all('/*', handler)


app.listen(3000);


// F U N C S -----------------------------------------
function handler(req, res) {
    var s = `
<<---------------------------
method: ${ st(req.method)}
headers: ${ st(req.headers)}
url:         ${ st(req.url)}
originalUrl: ${ st(req.originalUrl)}
path:        ${ st(req.path)}
query: ${ st(req.query)}
body: ${ st(req.body)}
--------------------------->>
`
    console.log(s)
    res.send(s)
    // res.end()
}

function st(params) {
    return JSON.stringify(params, null, '  ')
}

