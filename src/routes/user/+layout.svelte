<script>
    import { UserRole } from "$lib/models";
    import AuthRoute from "../../components/auth/AuthRoute.svelte";
</script>

<AuthRoute roles={[ UserRole.USER, UserRole.ADMIN ]}>
    <slot/>
</AuthRoute>