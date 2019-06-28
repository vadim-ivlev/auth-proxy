const process = require('process');
var app = require('express')();
var bodyParser = require('body-parser');

app.use(bodyParser.json()); // for parsing application/json
app.use(bodyParser.urlencoded({ extended: true })); // for parsing application/x-www-form-urlencoded

// app.post('/*', handler)
// app.get('/*', handler)
// app.options('/*', handler)
app.all('/*', handler)


var port = process.argv[2]
var appName = process.argv.slice(3).join(' ')


console.log(`"${appName}" started at http://localhost:${port}`)
app.listen(port);




// F U N C S -----------------------------------------
function handler(req, res) {
    var info = `

app: ${appName}    
---------------------------
method: ${ st(req.method)}
headers: ${ st(req.headers)}
url:         ${ st(req.url)}
originalUrl: ${ st(req.originalUrl)}
path:        ${ st(req.path)}
query: ${ st(req.query)}
body: ${ st(req.body)}
---------------------------
`


    

    res.send(`
    <!DOCTYPE html>
    <head>
    </head>
    <body>
        <h1>${appName}</h1>
        <pre>
            ${info}
        </pre>
    </body>
    </html>    
    `)

    console.log(info)

    // res.end()
}

function st(params) {
    return JSON.stringify(params, null, '  ')
}

