// M O D E L  ******************************************************************************************

var model = {
    origin: document.location.origin,
    _user: null,
    _app: null,
    _apps: null,
    _users: null,
    get logined(){
        return (getCookie('auth-proxy') != "")
    },
    set user(v) {
        this._user = v
        renderPage('user','#userPage')
    },
    get user() {
        return this._user
    },
    set app(v) {
        this._app = v
        renderPage('app','#appPage')
    },
    get app() {
        return this._app
    },
    set users(v) {
        this._users = v
        renderPage('users','.user-search-results')
    },
    get users() {
        return this._users
    },
    set apps(v) {
        this._apps = v
        renderPage('apps','.app-search-results')
    },
    get apps() {
        return this._apps
    },

}


// F U N C T I O N S  *********************************************************************************

function showPage(pageid, elemSelector){
    // $('.tab').css("border-top","1px solid transparent")
    // $('#'+pageid+'Tab').css("border-top","1px solid #9b4dca")

    $('.page').hide()
    $('#'+pageid+'Page').show()
    return false
}


// blinkStatus показывает исчезающее сообщение
function blinkStatus(selector, message) {
    console.log("blink:", message)
    let st = $(selector)
    st.text(message)
    st.fadeTo(0,1)
    st.fadeTo(2000, 0.001)
}


function renderTemplateFile(templateFile, data, targetSelector) {
    $.get(templateFile, function(template) {
        var rendered = Mustache.render(template, data);
        $(targetSelector).html(rendered);
    });
}


function renderPage(pageid, elemSelector) {
    renderTemplateFile('templates/mustache/'+pageid+'.html', model, elemSelector)
}


function renderMenu(){
    renderTemplateFile('templates/mustache/menu.html', model, '#menu')
}


function alertOnError(e, msg){
    alert(msg, e)
}

// function showResponseIn(elementID, showLoginPage=false) {
//     return function(response){
//         $(elementID).jsonViewer(response, {collapsed: true, rootCollapsable: false})
//         renderMenu()
//         if (showLoginPage) showPage(model.logined ?'logout': 'login')
//     }
// } 

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






/*
// R E Q U E S T S  *******************************************************

// function loginRestFormSubmit(event) {
//     if (event) event.preventDefault()
//     $("#resultLoginRest").html("")
//     $("#resultLogoutRest").html("")
//     $("#formLoginRest").ajaxSubmit({
//         url: "/login",
//         type: "POST",
//         success: showResponseIn('#resultLoginRest',true),
//         error: alertOnError
//     })
//     return false       
// }

// function logoutRestFormSubmit(event) {
//     if (event) event.preventDefault()
//     $("#resultLoginRest").html("")
//     $("#resultLogoutRest").html("")
//     $("#formLogoutRest").ajaxSubmit({
//         url: "/logout",
//         type: "GET",
//         success: showResponseIn('#resultLogoutRest',true),
//         error: alertOnError
//     })
//     return false
// }
*/

// L O G I N  *****************************************************************************

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
    $.ajax({ url: "/graphql", type: "POST", data: { query: query }, error: alertOnError,
        success: (res) => {
            showJSON(res,'#resultLoginGraphQL')
            renderMenu();
            if (res.errors){
                blinkStatus("#loginStatus", "Пароль или имя или емайл не подходят")
                return
            }
            formListAppSubmit();  
            showPage('apps') ;
        }   
    })
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
    $.ajax({ url: "/graphql", type: "POST", data: { query: query }, error: alertOnError,
        success: (res) => {
            showJSON(res,'#resultLogoutGraphQL')
            renderMenu();
            showPage('login') ;
        }
       
    })
    return false       
}

// U S E R S  *******************************************************************

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
    $.ajax({ url: "/graphql", type: "POST", data: { query: query }, error: alertOnError,
        success: (res) => {
            showJSON(res,'#resultListUser')
            model.users = res.data.list_user.list
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
            description
            email
            fullname
            password
            username
          }

        }
    `
    $.ajax({ url: "/graphql", type: "POST", data: { query: query }, error: alertOnError,
        success: (res) => {
            showJSON(res,'#resultUser')
            if (res.errors){
                blinkStatus("#userStatus", res.errors[0].message)
                return
            }
            model.user = res.data[userOperationName]
            setTimeout(()=>blinkStatus("#userStatus", userOperationName+" success" ), 100)
            getUser(username)
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
            description
            email
            fullname
            password
            username
          }
        
        list_app_user_role(
        username: "${username}"
        ) {
            app_description
            appname
            rolename
            username
          }
        
        }

    `

    $.ajax({ url: "/graphql", type: "POST", data: { query: query }, error: alertOnError,
        success: (res) => {
            showJSON(res,'#resultUser')
            if (res.errors){
                blinkStatus("#userStatus", res.errors[0].message)
                return
            }
            model.user = res.data.get_user
            model.user.apps = groupApps(res.data.list_app_user_role)
            setTimeout(()=>blinkStatus("#userStatus", "Значение в бд" ), 100)
        } 
    })
    return false       
}


function groupApps(list_app_user_role) {
    let gr = {}
    for (let aur of list_app_user_role ){
        gr[aur.appname] =[]
    }
    for (let aur of list_app_user_role ){
        gr[aur.appname].push(aur)
    }

    let arr = []

    for (let [key, value] of Object.entries(gr)) {
        console.log(`${key}: ${value}`);
        let rec = {}
        rec.appname =key
        rec.app_description = value[0].app_description
        rec.items = value
        arr.push(rec)
    }
    return arr
}



function deleteUser(username) {
    $("#resultUser").html("")
    var query =`
    mutation {
        delete_user(
        username: "${username}"
        ) {
            username
          }
        }
    `
    $.ajax({ url: "/graphql", type: "POST", data: { query: query }, error: alertOnError,
        success: (res) => {
            showJSON(res,'#resultUser')
            if (res.errors){
                blinkStatus("#userStatus", res.errors[0].message)
                return
            }
            model.user = null
            formListUserSubmit();  
            showPage('users') ;
         } 
    })
    return false       
}


// A P P S  *******************************************************************

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
              url
            }
          }
        }    `
    $.ajax({ url: "/graphql", type: "POST", data: { query: query }, error: alertOnError,
        success: (res) => {
            showJSON(res,'#resultListApp')
            model.apps = res.data.list_app.list
        } 
    
    })
    return false       
}



// A P P   U S E R   R O L E   **************************************************************************************************

function formListRoleSubmit(event) {
    if (event) event.preventDefault()
    // $("#resultListUser").html("")
    // let search = $("#formListUser input[name='search']").val()
    // var query =`
    // query {
    //     list_user(
    //     search: "${search}",
    //     order: "fullname ASC"
    //     ) {
    //         length
    //         list {
    //           description
    //           email
    //           fullname
    //           username
    //         }
    //       }
    //     }        
    //     `
    // $.ajax({ url: "/graphql", type: "POST", data: { query: query }, error: alertOnError,
    //     success: (res) => {
    //         showJSON(res,'#resultListUser')
    //         model.users = res.data.list_user.list
    //     } 
    
    // })
    return false       
}





// O N   P A G E   L O A D  ****************************************************************************************


renderMenu()
// renderPage('user')
// renderPage('app')

if (model.logined) {
    formListAppSubmit();  
    showPage('apps') ;
} else {
    renderPage('login','#loginPage');
    showPage('login')
}