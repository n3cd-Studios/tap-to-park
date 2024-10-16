<script>

    import AuthRoute from "../../../components/auth/AuthRoute.svelte";
    import { authStore } from "$lib/auth";

</script>

<AuthRoute pass={!$authStore.token} redirect="/register" signOut={false}>
    <slot />
</AuthRoute>