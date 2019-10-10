// M O D E L  ******************************************************************************************

var model = {
    // if debug == true logs go to console.
    debug: true,
    origin: document.location.origin,


    get logined(){
         return (this._loginedUser != null)
    },
 
    //---------------------------
    _loginedUser: null,
    set loginedUser(v) {
        this._loginedUser = v
        if (v) {
            document.getElementById("userTab").innerText = v.username
            // getAuthRoles(model.loginedUser.username)
        } else {
            document.getElementById("userTab").innerText = ""
        }
        refreshApp()
        
    },
    get loginedUser() {
        return this._loginedUser
    },

    //---------------------------
    _selfRegAllowed: false,
    set selfRegAllowed(v) {
        this._selfRegAllowed = v
        if (v) {
            showElements("#selfRegButton")
        } else {
            hideElements("#selfRegButton")
        }
    },
    get selfRegAllowed() {
        return this._selfRegAllowed
    },
    
    //---------------------------
    _captchaRequired: false,
    set captchaRequired(v) {
        this._captchaRequired = v
        if (v) {
            showElements("#captcha")
        } else {
            hideElements("#captcha")
        }
    },
    get captchaRequired() {
        return this._captchaRequired
    },

    //---------------------------
    _authRoles: null,
    set authRoles(v) {
        this._authRoles = v
        if (this.isAdmin){
            showElements('#usersTab')
            showElements('#rolesTab')
            showElements('#graphqlTest')
            showElements('#btnNewApp')
        } else {
            hideElements('#usersTab')
            hideElements('#rolesTab')
            hideElements('#graphqlTest')
            hideElements('#btnNewApp')

            // showPage('apps')
        }
        renderPage('apps','.app-search-results')
        
    },
    get authRoles() {
        return this._authRoles
    },
    get isAdmin() {
        if (!model.authRoles) return false
        return model.authRoles.some(e => e.rolename == "authadmin")
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
        this.all_app_options = createOptions(v, "appname", "description", "url")
        document.querySelector("#allAppsDataList").innerHTML = this.all_app_options
    },
    get allApps() {
        return this._allApps
    },

    //---------------------------
    _allUsers: null,
    set allUsers(v) {
        this._allUsers = v
        this.all_user_options = createOptions(v, "username", "fullname", "email")
        document.querySelector("#allUsersDataList").innerHTML = this.all_user_options
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

function createOptions(selectValues, keyProp, textProp1, textProp2) {
    var output = []
    selectValues && selectValues.forEach(function(value)
    {
      output.push(`<option value="${value[keyProp]}">${value[textProp1]} &nbsp;&nbsp;&nbsp; ${value[textProp2]?value[textProp2]:''}</option>`);
    })
    let optionText = output.join('')
    return optionText
}


function highlightTab(tabid) {
    removeClass('.tab', "underlined")
    var tabid0 = tabid.split("/")[0]
    addClass('#'+tabid0+'Tab', "underlined")   
}



function showPage(pageid, dontpush){
    //распарсить pageidExtended
    var a = pageid.split("/")
    var pageid0 = a[0]
    var id = a[1]

    highlightTab(pageid)
    
    hideElements('.page')
    showElements('#'+pageid0+'Page')


    // setting focus
    var text = document.querySelector('#'+pageid0+'Page input[type="text"]')
    if(text) 
        text.focus()


    if (!dontpush){
        if (!history.state || history.state.pageid != pageid ){
            history.pushState({pageid:pageid},pageid, "#"+pageid) 
        }
    }

    if (id) {
        if (pageid0 == "app"){
            getApp(id)
        } 
        if (pageid0 == "user"){
            getUser(id)
        }
    }


    return false
}



function renderTemplateFile(templateFile, data, targetSelector) {

    function onSuccess(template) {
        var rendered = Mustache.render(template, data)
        document.querySelector(targetSelector).innerHTML = rendered
    }    

    // $.get(templateFile, onSuccess);
    fetch(templateFile).then(x => x.text()).then(onSuccess)
}


function renderPage(pageid, elemSelector) {
    renderTemplateFile('templates/mustache/'+pageid+'.html', model, elemSelector)
}


function alertOnError(e, msg){
    alert(msg, e)
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


function searchApps() {
    if (document.querySelector('#chkLocalSearch').checked) {
        return delayFunc(searchAppsInModel, 100)
    } else {
        return delayFunc(formListAppSubmit)
    }
}

function searchAppsInModel() {
    if (!model.allApps) return
    var text = document.querySelector("#formListApp input[name='search']").value.trim().replace(' ','.*')
    var r = new RegExp(text, 'i')
    let found = model.allApps.filter((v)=>{
        var s = Object.values(v).join(' ')
        return r.test(s)
    })
    model.apps = found
    return false   
}


function sortAppsBy(prop) {
    if (!model._allApps) return false
    model._allApps.sort( (a,b) => (a[prop]+a.appname)>(b[prop]+b.appname)? 1: -1)
    model._apps.sort( (a,b) => (a[prop]+a.appname)>(b[prop]+b.appname)? 1: -1 )
    model.apps = model._apps
    return false
}



function searchUsers() {
    if (document.querySelector('#chkLocalSearch').checked) {
        return delayFunc(searchUsersInModel, 100)
    } else {
        return delayFunc(formListUserSubmit)
    }
}

function searchUsersInModel() {
    if (!model.allUsers) return
    var text = document.querySelector("#formListUser input[name='search']").value.trim().replace(' ','.*')
    var r = new RegExp(text, 'i')
    let found = model.allUsers.filter( (v) => r.test (Object.values(v).join(' ')) )
    model.users = found
    return false   
}

function sortUsersBy(prop) {
    if (!model._allUsers) return false
    model._allUsers.sort( (a,b) => (a[prop]+a.username)>(b[prop]+b.username)? 1: -1)
    model._users.sort( (a,b) => (a[prop]+a.username)>(b[prop]+b.username)? 1: -1 )
    model.users = model._users
    return false
}



// R E Q U E S T S  *******************************************************

function doGraphQLRequest(query, responseHandler, errorElementID) {
    // $.ajax({ url: "/graphql", type: "POST", data: { query: query }, error: alertOnError, 
    //     success: (res) => {
    //         model.debug && console.log(res)
    //         if (res.errors){
    //             model.debug && console.log(res.errors[0].message)
    //             if (errorElementID) {
    //                 document.getElementById(errorElementID).innerText = res.errors[0].message
    //             }
    //             return
    //         }
    //         responseHandler(res)
    //     } 
    // })

    fetch('/graphql', { 
        method: 'POST', 
        credentials: 'include', 
        body: JSON.stringify({ query: query, variables: {} }) 
    })
        .then(res => { 
          if (res.ok) return res.json();
          new Error(res)
        })
        .then((res) => {
            model.debug && console.log(res)
            if (res.errors){
                model.debug && console.log(res.errors[0].message)
                if (errorElementID) {
                    document.getElementById(errorElementID).innerText = res.errors[0].message
                }
                return
            }
            responseHandler(res)
        })
        .catch(console.error)    
    
}


// L O G I N  *****************************************************************************

function loginGraphQLFormSubmit(event) {
    if (event) event.preventDefault()
    
    let username = document.querySelector("#formLoginGraphQL input[name='username']").value
    let password = document.querySelector("#formLoginGraphQL input[name='password']").value
    let captcha =  document.querySelector("#formLoginGraphQL input[name='captcha']").value
    

    let query =`
    query {
        login(
        username: "${username}",
        password: "${password}",
        captcha: "${captcha}"
        )
        }    
    `
    function onSuccess(res){
        getLoginedUser()
    }   

    doGraphQLRequest(query, onSuccess, "loginError")
    clearLoginForm()
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

    function onSuccess(res){
        model.loginedUser = null
        refreshApp()
    }

    doGraphQLRequest(query, onSuccess)
    return false       
}



function isSelfRegAllowed(event) {
    if (event) event.preventDefault()
    var query =` query { is_selfreg_allowed }`

    function onSuccess(res){
        model.selfRegAllowed = res.data.is_selfreg_allowed
    }
       
    doGraphQLRequest(query, onSuccess)
    return false       
}

function isCaptchaRequired(event) {
    if (event) event.preventDefault()
    let username = document.querySelector("#formLoginGraphQL input[name='username']").value   
    var query =`  query { is_captcha_required(  username: "${username}" ) } `

    function onSuccess(res){
        model.captchaRequired = res.data.is_captcha_required
    }
       
    doGraphQLRequest(query, onSuccess)
    return false       
}



function generateNewPassword(event) {
    if (event) event.preventDefault()

    let username = document.querySelector("#formLoginGraphQL input[name='username']").value

    var query =`
    mutation {
        generate_password(
        username: "${username}"
        ) 
        }
    `
    function onSuccess(res){
        alert(res.data.generate_password)
        refreshApp()
    }   

    doGraphQLRequest(query, onSuccess)
    return false       
}



// L O G I N E D   U S E R   **********************************************************************************************************************


function getLoginedUser() {
    model.loginedUser = null
    
    var query =`
    query {
        get_logined_user {
            description
            email
            fullname
            username
            disabled
        }
    }
    `
    
    function onSuccess(res){
        model.loginedUser = res.data.get_logined_user
        getAuthRoles(model.loginedUser.username)
        getUser(model.loginedUser.username)
    } 
    
    doGraphQLRequest(query, onSuccess)
    return false       
}

function getAuthRoles(username) {
    model.authRoles = null

    var query =`
    query {
        list_app_user_role(
            appname: "auth",
            username: "${username}"
            ) {
                rolename
            }
        }
        `

    function onSuccess(res){
        model.authRoles = res.data.list_app_user_role
    } 

    doGraphQLRequest(query, onSuccess)
    return false       
}


// U S E R S  ***********************************************************************************************************************


function formListUserSubmit(event) {
    if (event) event.preventDefault()
    model.users = null
    let search = document.querySelector("#formListUser input[name='search']").value
    
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
              disabled
            }
          }
        }        
    `
    
    function onSuccess(res){
        model.users = res.data.list_user.list
    } 

    doGraphQLRequest(query, onSuccess)
    return false       
}



function formUserSubmit(event, userOperationName = 'create_user') {
    if (event) event.preventDefault()
    let username =      document.querySelector("#formUser input[name='username']").value
    let password =      document.querySelector("#formUser input[name='password']").value
    let email    =      document.querySelector("#formUser input[name='email']").value
    let fullname =      document.querySelector("#formUser input[name='fullname']").value
    let description =   document.querySelector("#formUser *[name='description']").value
    let disabled =      document.querySelector("#formUser input[name='disabled']").value
    
    var query =`
    mutation {
        ${userOperationName}(
        username: "${username}",
        password: "${password}",
        email: "${email}",
        fullname: "${fullname}",
        description: "${description}",
        disabled: ${disabled}
        ) {
            description
            email
            fullname
            username
            disabled
          }

        }
    `

    function onSuccess(res){
        alert(userOperationName+" success")
        refreshData()
        model.user = res.data[userOperationName]
        if (userOperationName == 'create_user' && !model.logined) {
            alert(`"${username}" is created.` )
            showPage('login',true)
        }
        getUser(username)
    }

    doGraphQLRequest(query, onSuccess)
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
            disabled
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

    function onSuccess(res){
        model.user = res.data.get_user
        model.user.apps = groupApps(res.data.list_app_user_role)
    } 

    doGraphQLRequest(query, onSuccess)
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
    function onSuccess(res){
        model.user = null
        refreshData()
        showPage('users')
    } 

    doGraphQLRequest(query, onSuccess)
    return false       
}


// A P P S  *******************************************************************

function formListAppSubmit(event) {
    if (event) event.preventDefault()
    model.apps = null
    let search = document.querySelector("#formListApp input[name='search']").value
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
              rebase
              public
            }
          }
        }    `

    function onSuccess(res){
        model.apps = res.data.list_app.list
    } 

    doGraphQLRequest(query, onSuccess)
    return false       
}



function formAppSubmit(event, appOperationName = 'create_app') {
    if (event) event.preventDefault()
    model.app = null
    let appname =     document.querySelector("#formApp input[name='appname']"    ).value
    let url =         document.querySelector("#formApp input[name='url']"        ).value
    let description = document.querySelector("#formApp input[name='description']").value
    let rebase =      document.querySelector("#formApp input[name='rebase']"     ).value
    let public =      document.querySelector("#formApp input[name='public']"     ).value
    
    var query =`
    mutation {
        ${appOperationName}(
        appname: "${appname}",
        url: "${url}",
        description: "${description}",
        rebase: "${rebase}",
        public: "${public}"
        ) {
            description
            appname
            url
            rebase
            public
          }

        }
    `
    function onSuccess(res){
        alert(appOperationName+" success")
        model.app = res.data[appOperationName]
        refreshData()
        getApp(appname)
    } 

    doGraphQLRequest(query, onSuccess)
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
            rebase
            public
          }
        
        list_app_user_role(
        appname: "${appname}"
        ) {
            appname
            rolename
            user_fullname
            user_disabled
            username
          }
        
        }

    `

    function onSuccess(res){
        model.app = res.data.get_app
        model.app.users = groupUsers(res.data.list_app_user_role)
    } 

    doGraphQLRequest(query, onSuccess)
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
    function onSuccess(res){
        model.app = null
        refreshData()
        showPage('apps') ;
    } 

    doGraphQLRequest(query, onSuccess)
    return false       
}




// A P P   U S E R   R O L E   **************************************************************************************************

function getAllApps(event) {
    if (event) event.preventDefault()
    model.allApps = null
    var query =`
    query {
        list_app(
        order: "appname ASC"
        ) {
            length
            list {
              appname
              description
              url
              rebase
              public 
            }
          }
        }    `

    function onSuccess(res){
        model.allApps = res.data.list_app.list
    } 

    doGraphQLRequest(query, onSuccess)
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
              disabled
            }
          }
        }    `

    function onSuccess(res){
        model.allUsers = res.data.list_user.list
    } 

    doGraphQLRequest(query, onSuccess)
    return false       
}

function formListRoleSubmit(event) {
    if (event && event.preventDefault ) event.preventDefault()
    model.app_user_roles = null
    let appname =  document.querySelector("#allApps").value
    let username = document.querySelector("#allUsers").value
    if (!appname || !username) 
        return

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

    function onSuccess(res){
        model.app_user_roles = res.data.list_app_user_role
    } 

    doGraphQLRequest(query, onSuccess)
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
            appname
            username
          }
        }
    `
    function onSuccess(res){
        if (onsuccess) onsuccess()
    }

    doGraphQLRequest(query, onSuccess)
    return false       
}

// works when input values on roles page change
function refreshRoles() {
    let allUsersValue = document.getElementById("allUsers").value
    if (allUsersValue) {
        let ui = document.querySelector(`#allUsersDataList>option[value='${allUsersValue}']`).innerText
        document.getElementById('userInfo').innerText = ui   
    }

    let allAppsValue = document.getElementById("allApps").value
    if (allAppsValue) {
        let ai = document.querySelector(`#allAppsDataList>option[value='${allAppsValue}']`).innerText
        document.getElementById('appInfo').innerText = ai
    }

    if (allUsersValue && allAppsValue) 
        formListRoleSubmit() 
}


function filterRows(selector, value ){
    var v = value.toLowerCase()
    var rows = document.querySelectorAll(selector)
    rows.forEach(e => {
        var txt = e.innerText.toLowerCase()
        if (txt.indexOf(v) == -1) {
            e.classList.add("hidden")
        } else {
            e.classList.remove("hidden")
        }
    });
}

function hideElements(selector) {
    document.querySelectorAll(selector).forEach(e => e.classList.add("hidden"));
}


function showElements(selector) {
    document.querySelectorAll(selector).forEach(e => e.classList.remove("hidden"));
}

function addClass(selector, classname) {
    document.querySelectorAll(selector).forEach(e => e.classList.add(classname));
}

function removeClass(selector, classname) {
    document.querySelectorAll(selector).forEach(e => e.classList.remove(classname));
}


function hidePassword() {
    document.querySelector("#formUser input[name='password']").setAttribute('type','password')
}

function showPassword() {
    document.querySelector("#formUser input[name='password']").setAttribute('type','text')
}

function generatePassword() {
    function newPassword (n) {
        let pickSymbol =(s) => s[Math.floor(Math.random()*s.length)]
        var symbolSets =["bcdfghjklmnpqrstvwxz","aeiou"] 
        var password = ''
        for (let i=0; i<n; i++){
            password += pickSymbol(symbolSets[i%symbolSets.length])
        }
        return password
    }

    document.querySelector("#formUser input[name='password']").value = newPassword(9)
}


function getNewCaptcha(event) {
    let uri = "captcha?"+ new Date().getTime()
    document.getElementById("captchaImg").src = uri
    return false
}


function clearLoginForm() {
    document.querySelector("#formLoginGraphQL input[name='username']").value = ""
    document.querySelector("#formLoginGraphQL input[name='password']").value = ""
    document.querySelector("#formLoginGraphQL input[name='captcha']").value = ""
    document.getElementById("loginError").innerText = ""
}


function logout() {
    logoutGraphQLFormSubmit()
    clearLoginForm()
    showPage('login',true)
    isSelfRegAllowed()
    model.captchaRequired = false
    isCaptchaRequired()
    return false
}


// O N   P A G E   L O A D  ****************************************************************************************

function refreshData() {
    if (model.logined) {
        // getLoginedUser()
        getAllApps()
        getAllUsers()
        formListAppSubmit()
        formListUserSubmit()  
    } else {
        // nullify model's inner props
        for (const k of Object.keys(model)) {
            if (k.startsWith('_')) {
                model[k] = null
                // console.log(`model.${k} = null`)
            }
        }
        isSelfRegAllowed()
   }    
}

function getLandingPageid(){
    var p = location.hash.slice(1)
    return p ? p : 'apps'
}

function refreshApp(params) {
    refreshData()

    if (model.logined) {
        let page = getLandingPageid()
        showPage(page) 
        showElements('#menu')     
    } else {
        showPage('login',true)
        hideElements('#menu')
    }  
}


window.onhashchange = function(event) {
    console.log("onhashchange", event)
    var newpage = event.newURL.split('#')[1]
    if (newpage) 
        showPage(newpage)
}

// window.onpopstate = function(event) {
//     console.log( "event.state: " + JSON.stringify(event.state));
// }


model.captchaRequired = false
getLoginedUser()
refreshApp()