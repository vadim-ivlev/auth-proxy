// M O D E L  ******************************************************************************************

var model = {
    origin: document.location.origin,

    get logined(){
        return (getCookie('auth-proxy') != "")
    },
 
    //---------------------------
    _loginedUser: null,
    set loginedUser(v) {
        this._loginedUser = v
        $("#loginedUser").text( v? v.username: '')
    },
    get loginedUser() {
        return this._loginedUser
    },
    
    
    //---------------------------
    _user: null,
    set user(v) {
        this._user = v
        renderPage('user','#userPage')
    },
    get user() {
        return this._user
    },
    
    //---------------------------
    _app: null,
    set app(v) {
        this._app = v
        renderPage('app','#appPage')
    },
    get app() {
        return this._app
    },
    
    //---------------------------
    _users: null,
    set users(v) {
        this._users = v
        renderPage('users','.user-search-results')
    },
    get users() {
        return this._users
    },
    
    //---------------------------
    _apps: null,
    set apps(v) {
        this._apps = v
        renderPage('apps','.app-search-results')
    },
    get apps() {
        return this._apps
    },

    //---------------------------
    all_app_options: null,
    all_user_options:null,
    //---------------------------
    _allApps: null,
    set allApps(v) {
        this._allApps = v
        this.all_app_options = createOptions(v, "appname", "description")
        $("#allApps").html(this.all_app_options)
        $("#allApps").val("app1")
    },
    get allApps() {
        return this._allApps
    },

    //---------------------------
    _allUsers: null,
    set allUsers(v) {
        this._allUsers = v
        this.all_user_options = createOptions(v, "username", "fullname")
        $("#allUsers").html(this.all_user_options)
        $("#allUsers").val("vadim")
    },
    get allUsers() {
        return this._allUsers
    },
    
    //---------------------------
    _app_user_roles: null,
    set app_user_roles(v) {
        this._app_user_roles = v
        renderPage('roles','.app-user-roles-results')
    },
    get app_user_roles() {
        return this._app_user_roles
    },

    
}



// F U N C T I O N S  *********************************************************************************

function createOptions(selectValues, keyProp, textProp) {
    var output = []
    $.each(selectValues, function(key, value)
    {
      output.push('<option value="'+ value[keyProp] +'">'+ value[textProp] +'</option>');
    })
    let optionText = output.join('')
    return optionText
}


function highlightTab(tabid) {
    $('.tab').css("border-bottom-color","transparent")
    $('#'+tabid+'Tab').css("border-bottom-color","#9b4dca")   
}



function showPage(pageid, dontpush){
    highlightTab(pageid)
    
    $('.page').hide()
    $('#'+pageid+'Page').show()

    // setting focus
    // var text = $('#'+pageid+'Page input[type="text"]')[0]
    // if(text) 
    //     text.focus()

    if (!dontpush){
        if (!history.state || history.state.pageid != pageid ){
            history.pushState({pageid:pageid},pageid, "#"+pageid) 
            console.log("push", pageid)   
        }
    }
    return false
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


function alertOnError(e, msg){
    alert(msg, e)
}


function showResponse(response) {
    if (! document.querySelector('#chkShowResponses').checked) return false;
    return showJSON(response)
}


function showJSON(model) {
    $('#jsonViewerView').html("")
    $('#jsonViewerView').jsonViewer(model, {collapsed: true, rootCollapsable: false}) 
    $('#jsonViewer').show()
    return false
}


// blinkStatus shows fading message
function blinkStatus(message) {
    console.log("blink:", message)
    let st = $("#blinkMessage")
    st.text(message)
    // st.show()
    st.fadeTo(0,1)
    st.fadeTo(2000, 0.0)
    // st.hide(2000)
}

function showModel() {
    $('#jsonViewerView').html("")
    $('#jsonViewerView').jsonViewer(model, {collapsed: true, rootCollapsable: false}) 
    $('#jsonViewer').show()
   return false
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
function delayFunc(f, delay=500) {
   clearTimeout(delayTimeout) 
   delayTimeout = setTimeout(f, delay)  
   return false 
}




/*
// R E Q U E S T S  *******************************************************

// function loginRestFormSubmit(event) {
//     if (event) event.preventDefault()
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
            showResponse(res)
            if (res.errors){
                blinkStatus("Пароль или имя или емайл не подходят")
                return
            }
            refreshApp()
        }   
    })
    return false       
}



function logoutGraphQLFormSubmit(event) {
    if (event) event.preventDefault()
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
            showResponse(res)
            if (res.errors){
                blinkStatus("Пароль или имя или емайл не подходят")
                return
            }
            refreshApp()
        }
       
    })
    return false       
}

// U S E R S  *******************************************************************


function getLoginedUser() {
    model.loginedUser = null
    var query =`
    query {
        get_logined_user {
            description
            email
            fullname
            username
          }
        }
    `

    $.ajax({ url: "/graphql", type: "POST", data: { query: query }, error: alertOnError,
        success: (res) => {
            showResponse(res)
            if (res.errors){
                blinkStatus( res.errors[0].message)
                return
            }
            model.loginedUser = res.data.get_logined_user
        } 
    })
    return false       
}




function formListUserSubmit(event) {
    if (event) event.preventDefault()
    model.users = null
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
            showResponse(res)
            if (res.errors){
                blinkStatus( res.errors[0].message)
                return
            }
            model.users = res.data.list_user.list
        } 
    
    })
    return false       
}



function formUserSubmit(event, userOperationName = 'create_user') {
    if (event) event.preventDefault()
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
            username
          }

        }
    `
    $.ajax({ url: "/graphql", type: "POST", data: { query: query }, error: alertOnError,
        success: (res) => {
            showResponse(res)
            if (res.errors){
                blinkStatus( res.errors[0].message)
                return
            }
            blinkStatus( userOperationName+" success" )
            refreshData()
            model.user = res.data[userOperationName]
            getUser(username)
        } 
    })
    return false       
}



function getUser(username) {
    model.user = null
    var query =`
    query {
        get_user(
        username: "${username}"
        ) {
            description
            email
            fullname
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
            showResponse(res)
            if (!res.data.get_user){
                blinkStatus( res.errors[0].message)
                return
            }
            model.user = res.data.get_user
            model.user.apps = groupApps(res.data.list_app_user_role)
        } 
    })
    return false       
}


function groupApps(list_app_user_role) {
    if (!list_app_user_role) return []
    
    let gr = {}
    for (let aur of list_app_user_role ){
        if (!gr[aur.appname]) gr[aur.appname] =[]
        gr[aur.appname].push(aur)
    }

    let arr = []

    for (let [key, value] of Object.entries(gr)) {
        // console.log(`${key}: ${value}`);
        let rec = {}
        rec.appname =key
        rec.app_description = value[0].app_description
        rec.items = value
        arr.push(rec)
    }
    return arr
}



function deleteUser(username) {
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
            showResponse(res)
            if (res.errors){
                blinkStatus( res.errors[0].message)
                return
            }
            model.user = null
            refreshData()
            showPage('users') ;
         } 
    })
    return false       
}


// A P P S  *******************************************************************

function formListAppSubmit(event) {
    if (event) event.preventDefault()
    model.apps = null
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
            showResponse(res)
            if (res.errors){
                blinkStatus(res.errors[0].message)
                return
            }
            model.apps = res.data.list_app.list
        } 
    
    })
    return false       
}



function formAppSubmit(event, appOperationName = 'create_app') {
    if (event) event.preventDefault()
    model.app = null
    let appname = $("#formApp input[name='appname']").val()
    let url = $("#formApp input[name='url']").val()
    let description = $("#formApp input[name='description']").val()
    var query =`
    mutation {
        ${appOperationName}(
        appname: "${appname}",
        url: "${url}",
        description: "${description}"
        ) {
            description
            appname
            url
          }

        }
    `
    $.ajax({ url: "/graphql", type: "POST", data: { query: query }, error: alertOnError,
        success: (res) => {
            showResponse(res)
            if (res.errors){
                blinkStatus( res.errors[0].message)
                return
            }
            blinkStatus( appOperationName+" success" )
            model.app = res.data[appOperationName]
            refreshData()
            getApp(appname)
        } 
    })
    return false       
}



function getApp(appname) {
    model.app = null
    var query =`
    query {
        get_app(
        appname: "${appname}"
        ) {
            description
            appname
            url
          }
        
        list_app_user_role(
        appname: "${appname}"
        ) {
            appname
            rolename
            user_fullname
            username
          }
        
        }

    `

    $.ajax({ url: "/graphql", type: "POST", data: { query: query }, error: alertOnError,
        success: (res) => {
            showResponse(res)
            if (res.errors){
                blinkStatus( res.errors[0].message)
                return
            }
            model.app = res.data.get_app
            model.app.users = groupUsers(res.data.list_app_user_role)
        } 
    })
    return false       
}


function groupUsers(list_app_user_role) {
    let gr = {}
    for (let aur of list_app_user_role ){
        if (!gr[aur.username]) gr[aur.username] =[]
        gr[aur.username].push(aur)
    }

    let arr = []

    for (let [key, value] of Object.entries(gr)) {
        let rec = {}
        rec.username =key
        rec.user_fullname = value[0].user_fullname
        rec.items = value
        arr.push(rec)
    }
    return arr
}



function deleteApp(appname) {
    var query =`
    mutation {
        delete_app(
        appname: "${appname}"
        ) {
            appname
          }
        }
    `
    $.ajax({ url: "/graphql", type: "POST", data: { query: query }, error: alertOnError,
        success: (res) => {
            showResponse(res)
            if (res.errors){
                blinkStatus( res.errors[0].message)
                return
            }
            model.app = null
            refreshData()
            showPage('apps') ;
         } 
    })
    return false       
}




// A P P   U S E R   R O L E   **************************************************************************************************

function getAllApps(event) {
    if (event) event.preventDefault()
    model.allApps = null
    var query =`
    query {
        list_app(
        order: "description ASC"
        ) {
            length
            list {
              appname
              description
            }
          }
        }    `
    $.ajax({ url: "/graphql", type: "POST", data: { query: query }, error: alertOnError,
        success: (res) => {
            if (res.errors){
                blinkStatus( res.errors[0].message)
                return
            }
            model.allApps = res.data.list_app.list
        } 
    
    })
    return false       
}

function getAllUsers(event) {
    if (event) event.preventDefault()
    model.allUsers = null
    var query =`
    query {
        list_user(
        order: "fullname ASC"
        ) {
            length
            list {
              username
              fullname
              email
              description
            }
          }
        }    `
    $.ajax({ url: "/graphql", type: "POST", data: { query: query }, error: alertOnError,
        success: (res) => {
            if (res.errors){
                blinkStatus( res.errors[0].message)
                return
            }
            model.allUsers = res.data.list_user.list
        } 
    
    })
    return false       
}

function formListRoleSubmit(event) {
    if (event && event.preventDefault ) event.preventDefault()
    model.app_user_roles = null
    let appname = $("#formListRole select[name='appname']").val()
    let username = $("#formListRole select[name='username']").val()
    if (!appname || !username) return

    var query =`
    query {
        list_app_user_role(
        appname: "${appname}",
        username: "${username}"
        ) {
            rolename
          }
        }        
        `
    $.ajax({ url: "/graphql", type: "POST", data: { query: query }, error: alertOnError,
        success: (res) => {
            showResponse(res)
            if (res.errors){
                blinkStatus( res.errors[0].message)
                return
            }
            model.app_user_roles = res.data.list_app_user_role
        } 
    
    })
    return false       
}


function modifyRole(action,appname,username,rolename, onsuccess ) {
    if (!appname || !username || !rolename) return

    var query =`
    mutation {
        ${action}_app_user_role(
        appname: "${appname}",
        username: "${username}",
        rolename: "${rolename}"
        ) {
            rolename
          }
        }
    `
    $.ajax({ url: "/graphql", type: "POST", data: { query: query }, error: alertOnError, 
    success: (res) => {
        if (res.errors){
            blinkStatus( res.errors[0].message)
            return
        }
        if (onsuccess) 
            onsuccess()
    }
    })
    return false       
}


// O N   P A G E   L O A D  ****************************************************************************************

function refreshData() {
    if (model.logined) {
        getLoginedUser()
        getAllApps()
        getAllUsers()
        formListAppSubmit()
        formListUserSubmit()  
        delayFunc(formListRoleSubmit)  
    } else {
        // nullify model's inner props
        for (const k of Object.keys(model)) {
            if (k.startsWith('_')) {
                model[k] = null
                console.log(`model.${k} = null`)
            }

        }
   }    
}

function getLandingPageid(){
    var p = location.hash.slice(1)
    return p ? p : 'apps'
}

function refreshApp(params) {
    refreshData()

    if (model.logined) {
        showPage(getLandingPageid()) 
        $('#loginButton').hide()
        $('#loginedArea').show()
        $('#menu').show()       
    } else {
        showPage('login',true)
        $('#loginButton').show()
        $('#loginedArea').hide()
        $('#menu').hide()
    }  
}


window.onpopstate = function(event) {

    // console.log( "event.state: " + JSON.stringify(event.state));
    if (event.state) {
        let nextPageid = event.state.pageid
        if ($('#userPage').is(':visible')){
            showPage('users', true)
            return
        }
        if ($('#appPage').is(':visible')){
            showPage('apps', true)
            return
        }
        showPage(event.state.pageid, true)
    }
}

refreshApp()