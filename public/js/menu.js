const router = new VueRouter( {
    mode: 'history',
    routes: [
        {
            path: '/config/index',
            name: 'index',
        },
        {
            path: '/config/project',
            name: 'project',
        },
        {
            path: '/config/callback-log',
            name: 'callbackLog',
        },
        {
            path: '/config/add-config',
            name: 'addConfig',
        },
    ]
})
