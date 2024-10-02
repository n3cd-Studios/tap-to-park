<script>

    import AuthRoute from "../../../components/AuthRoute.svelte";
    import { authStore } from "$lib/auth";

</script>

<AuthRoute pass={!$authStore.token} redirect="/admin" signOut={false}>
    <slot />
</AuthRoute>
