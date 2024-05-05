const userRoutes = [
    {
        path: '/register',
        name: 'register',
        component: () => import('@/views/user/RegisterView.vue'),
    },
    {
        path: '/login',
        name: 'login',
        component: () => import('@/views/user/LoginView.vue'),
    },
    {
        path: '/profile',
        name: 'profile',
        meta: {
            auth: true,
        },
        component: () => import('@/views/user/layout_Profile.vue'),
    },
    {
        path: '/editProfile',
        name: 'editProfile',
        component: () => import('@/views/user/editProfile.vue')
    },
];

export default userRoutes;
