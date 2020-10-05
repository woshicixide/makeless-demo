import Makeless from '@makeless/makeless-ui/src/makeless';

declare module 'vue/types/vue' {
    interface Vue {
        $makeless: Makeless;
        userLoaded: boolean;

        onUserLoaded(): void
    }
}
