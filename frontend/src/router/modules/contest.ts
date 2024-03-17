const Layout = () => import("@/layout/index.vue");

export default {
  path: "/",
  name: "Contest",
  component: Layout,
  redirect: "/contest",
  meta: {
    title: "竞赛",
    icon: "academicons:acm",
    rank: 0
  },
  children: [
    {
      path: "/contest",
      name: "",
      component: () => import("@/views/contest/index.vue"),
      meta: {
        title: "竞赛",
      }
    }
  ]
} satisfies RouteConfigsTable;
