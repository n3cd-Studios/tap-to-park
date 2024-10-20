<script>

    import AuthRoute from "../../../components/auth/AuthRoute.svelte";

</script>

<AuthRoute authorized={false} signOut={false} redirect="/admin">
    <slot />
</AuthRoute>