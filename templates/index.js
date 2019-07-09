// M O D E L  ******************************************************************************************

var model = {
    user: null,
    app: null,
    apps: [],
    users: [],
    get logined(){
        return getCookie('auth-proxy') != ""
    },
}


// F U N C T I O N S  *********************************************************************************

function showPage(pageid, elemSelector){
    renderPage(pageid, elemSelector)
    $('.page').hide()
    $('#'+pageid+'Page').show()
    return false
}

// blinkStatus показывает исчезающее сообщение
function blinkStatus(selector, message) {
    let status = $(selector)
    status.text(message)
    status.show()
    status.hide(1000)
}



function renderTemplateFile(templateFile, data, targetSelector) {
    $.get(templateFile, function(template) {
        var rendered = Mustache.render(template, data);
        $(targetSelector).html(rendered);
      });
}



function renderPage(pageid, elemSelector) {
    renderTemplateFile('templates/mustache/'+pageid+'.html', model, 
    (elemSelector ? elemSelector : '#'+pageid+'Page')
    )
}

function renderMenu(){
    renderTemplateFile('templates/mustache/menu.html', model, '#menu')
}




function onError(e, msg){
    alert(msg, e)
}

function showResponseIn(elementID, showLoginPage=false) {
    return function(response){
        $(elementID).jsonViewer(response, {collapsed: true, rootCollapsable: false})
        renderMenu()
        if (showLoginPage) showPage(model.logined ?'logout': 'login')
    }
} 

function showJSON(response, elementID) {
     $(elementID).jsonViewer(response, {collapsed: true, rootCollapsable: false}) 
}

function getCookie(cname) {
    var name = cname + "=";
    var decodedCookie = decodeURIComponent(document.cookie);
    var ca = decodedCookie.split(';');
    for(var i = 0; i <ca.length; i++) {
      var c = ca[i];
      while (c.charAt(0) == ' ') {
        c = c.substring(1);
      }
      if (c.indexOf(name) == 0) {
        return c.substring(name.length, c.length);
      }
    }
    return "";
}


var delayTimeout
function delayFunc(f) {
   clearTimeout(delayTimeout) 
   delayTimeout = setTimeout(f, 500)  
   return false 
}


// R E N D E R I N G  *********************************************************







// R E Q U E S T S  *******************************************************

function loginRestFormSubmit(event) {
    if (event) event.preventDefault()
    $("#resultLoginRest").html("")
    $("#resultLogoutRest").html("")
    $("#formLoginRest").ajaxSubmit({
        url: "/login",
        type: "POST",
        success: showResponseIn('#resultLoginRest',true),
        error: onError
    })
    return false       
}

function logoutRestFormSubmit(event) {
    if (event) event.preventDefault()
    $("#resultLoginRest").html("")
    $("#resultLogoutRest").html("")
    $("#formLogoutRest").ajaxSubmit({
        url: "/logout",
        type: "GET",
        success: showResponseIn('#resultLogoutRest',true),
        error: onError
    })
    return false
}


function loginGraphQLFormSubmit(event) {
    if (event) event.preventDefault()
    $("#resultLoginGraphQL").html("")
    $("#resultLogoutGraphQL").html("")

    let username = $("#formLoginGraphQL input[name='username']").val()
    let password = $("#formLoginGraphQL input[name='password']").val()

    var query =`
    query {
        login(
        username: "${username}",
        password: "${password}"
        )
        }    
    `
    $.ajax({ url: "/graphql", type: "POST", data: { query: query }, error: onError,
        success: showResponseIn('#resultLoginGraphQL',true) })
    return false       
}


function logoutGraphQLFormSubmit(event) {
    if (event) event.preventDefault()
    $("#resultLoginGraphQL").html("")
    $("#resultLogoutGraphQL").html("")
    var query =`
    query {
        logout {
            message
            username
          }
        }
    `
    $.ajax({ url: "/graphql", type: "POST", data: { query: query }, error: onError,
        success: showResponseIn('#resultLogoutGraphQL',true) })
    return false       
}

function formListAppSubmit(event) {
    if (event) event.preventDefault()
    $("#resultListApp").html("")
    let search = $("#formListApp input[name='search']").val()
    var query =`
    query {
        list_app(
        search: "${search}",
        order: "description ASC"
        ) {
            length
            list {
              appname
              description
            }
          }
        }    `
    $.ajax({ url: "/graphql", type: "POST", data: { query: query }, error: onError,
        success: (res) => {
            showJSON(res,'#resultListApp')
            model.apps = res.data.list_app.list
            renderPage('apps','.app-search-results')
        } 
    
    })
    return false       
}


function formListUserSubmit(event) {
    if (event) event.preventDefault()
    $("#resultListUser").html("")
    let search = $("#formListUser input[name='search']").val()
    var query =`
    query {
        list_user(
        search: "${search}",
        order: "fullname ASC"
        ) {
            length
            list {
              description
              email
              fullname
              username
            }
          }
        }        
        `
    $.ajax({ url: "/graphql", type: "POST", data: { query: query }, error: onError,
        success: (res) => {
            showJSON(res,'#resultListUser')
            model.users = res.data.list_user.list
            renderPage('users','.user-search-results')
        } 
    
    })
    return false       
}


function formUserSubmit(event, userOperationName = 'create_user') {
    if (event) event.preventDefault()
    $("#resultUser").html("")
    let username = $("#formUser input[name='username']").val()
    let password = $("#formUser input[name='password']").val()
    let email = $("#formUser input[name='email']").val()
    let fullname = $("#formUser input[name='fullname']").val()
    let description = $("#formUser input[name='description']").val()
    var query =`
    mutation {
        ${userOperationName}(
        username: "${username}",
        password: "${password}",
        email: "${email}",
        fullname: "${fullname}",
        description: "${description}"
        ) {
            apps {
                appname
            }
            description
            email
            fullname
            password
            username
          }
        }
    `
    $.ajax({ url: "/graphql", type: "POST", data: { query: query }, error: onError,
        success: (res) => {
            showJSON(res,'#resultUser')
            model.user = res.data[userOperationName]
            renderPage('user')
            blinkStatus("#formUser .status", userOperationName+" success" )
        } 
    })
    return false       
}

function getUser(username) {
    $("#resultUser").html("")
    var query =`
    query {
        get_user(
        username: "${username}"
        ) {
            apps {
              appname
            }
            description
            email
            fullname
            password
            username
          }
        }    
    `
    $.ajax({ url: "/graphql", type: "POST", data: { query: query }, error: onError,
        success: (res) => {
            showJSON(res,'#resultUser')
            model.user = res.data.get_user
            renderPage('user')
            blinkStatus("#formUser .status", "Значение в бд" )
        } 
    })
    return false       
}






// O N   P A G E   L O A D  ****************************************************************************************


renderMenu()
showPage(model.logined ?'logout': 'login')