const blogRoutes = [
    {
        path: '/blog',
        component: () => import('@/views/blog/blog.vue'),
        children: [
            {
                path: '',
                component: () => import('@/views/blog/article.vue'),
            },
            {
                path: '/blog/add',
                component: () => import('@/views/blog/article_add.vue'),
            },
            {
                path: '/blog/manage',
                component: () => import('@/views/blog/article_manage.vue'),
            }
        ],
    },
];

export default blogRoutes;
