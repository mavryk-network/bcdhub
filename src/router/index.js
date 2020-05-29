import Vue from 'vue'
import Router from 'vue-router'
import VueRouter from 'vue-router'

import Home from '@/views/home/Home.vue'
import Welcome from '@/views/home/Welcome.vue'

import ExtendedSearch from '@/views/extended_search/ExtendedSearch.vue'

import Stats from '@/views/stats/Stats.vue'
import NetworkTab from '@/views/stats/NetworkTab.vue'

import Diff from '@/views/diff/Diff.vue'

import Contract from '@/views/contract/Contract.vue'
import OperationsTab from '@/views/contract/OperationsTab.vue'
import CodeTab from '@/views/contract/CodeTab.vue'
import InteractTab from '@/views/contract/InteractTab.vue'
import StorageTab from '@/views/contract/StorageTab.vue'
import LogTab from '@/views/contract/LogTab.vue'

import OperationGroup from '@/views/opg/OperationGroup.vue'
import OpgContents from '@/views/opg/ContentsTab.vue'

import BigMap from '@/views/big_map/BigMap.vue'
import BigMapKeys from '@/views/big_map/KeysTab.vue'
import BigMapHistory from '@/views/big_map/HistoryTab.vue'

import Dashboard from '@/views/dashboard/Dashboard.vue'
import EventsTab from '@/views/dashboard/EventsTab.vue'


Vue.use(VueRouter)

const router = new Router({
    linkActiveClass: '',
    linkExactActiveClass: '',
    mode: 'history',
    routes: [
        {
            path: '/',
            components: {
                default: Home
            },
            name: 'home'
        },
        {
            path: '/welcome',
            components: {
                default: Welcome,
            },
            name: 'welcome',
            props: { default: true }
        },
        {
            path: '/search',
            components: {
                default: ExtendedSearch
            },
            name: 'search',
            props: { default: true },
        },
        {
            path: '/stats',
            components: {
                default: Stats,
            },
            name: 'stats',
            props: { default: true },
            children: [
                {
                    path: ':network',
                    name: 'network_stats',
                    component: NetworkTab,
                    props: true
                }
            ]
        },      
        {
            path: '/diff',
            components: {
                default: Diff,
            },
            name: 'diff',
            props: { default: true },
        },
        // {
        //     path: '/projects',
        //     components: {
        //         default: Projects,
        //         nav: Nav
        //     },
        //     name: 'projects'
        // },      
        { // backward compatibility
            path: '/:network(main|babylon|zero|carthage)/:address(KT[0-9A-z]{34})',
            children: [
                {
                    path: '',
                    redirect: to => {
                        const { params } = to
                        return `/${params.network}net/${params.address}`
                    }
                },
                {
                    path: 'operations',
                    redirect: to => {
                        const { params } = to
                        return `/${params.network}net/${params.address}/operations`
                    }
                },
                {
                    path: 'script',
                    redirect: to => {
                        const { params } = to
                        return `/${params.network}net/${params.address}/code`
                    }
                },
                {
                    path: 'state',
                    redirect: to => {
                        const { params } = to
                        return `/${params.network}net/${params.address}/storage`
                    }
                },
            ]
        },
        {
            path: '/:network/:address(KT[0-9A-z]{34})',
            components: {
                default: Contract,
            },
            props: { default: true },
            children: [
                {
                    path: '',
                    name: 'contract',
                    redirect: 'operations'
                },
                {
                    path: 'operations',
                    name: 'operations',
                    component: OperationsTab,
                    props: true
                },
                {
                    path: 'code',
                    name: 'code',
                    component: CodeTab,
                    props: true
                },
                {
                    path: 'interact',
                    name: 'interact',
                    component: InteractTab,
                    props: true
                },
                {
                    path: 'storage',
                    name: 'storage',
                    component: StorageTab,
                    props: true
                },
                {
                    path: 'log',
                    name: 'log',
                    component: LogTab,
                    props: true
                }
            ]
        },
        {
            path: '/dashboard',
            components: {
                default: Dashboard
            },
            props: { default: true },
            children: [
                {
                    path: '',
                    name: 'dashboard',
                    redirect: 'events'
                },
                {
                    path: 'events',
                    name: 'events',
                    component: EventsTab,
                    props: true
                }
            ]
        },
        {
            path: '/:network/big_map/:ptr(\\d+)',
            components: {
                default: BigMap
            },
            props: { default: true },
            children: [
                {
                    path: '',
                    name: 'big_map',
                    redirect: 'keys'
                },
                {
                    path: 'keys',
                    name: 'big_map_keys',
                    component: BigMapKeys,
                    props: true
                },
                {
                    path: ':keyhash',
                    name: 'big_map_history',
                    component: BigMapHistory,
                    props: true
                }
            ]
        },
        {
            path: '/:network/opg/:hash(o[0-9A-z]{50})',
            alias: '/:network(main|babylon|zero|carthage)/:hash(o[0-9A-z]{50})',
            components: {
                default: OperationGroup
            },
            props: { default: true },
            children: [
                {
                    path: '',
                    name: 'operation_group',
                    redirect: 'contents'
                },
                {
                    path: 'contents',
                    name: 'opg_contents',
                    component: OpgContents,
                    props: true
                }
            ]
        }
    ]
});

export default router;