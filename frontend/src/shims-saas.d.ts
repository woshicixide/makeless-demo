import Saas from '@go-saas/go-saas-ui/src/saas';

declare module 'vue/types/vue' {
    interface Vue {
        $saas: Saas;
    }
}
