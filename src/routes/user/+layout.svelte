<script>
    import { UserRole } from "$lib/models";
    import AuthRoute from "../../components/auth/AuthRoute.svelte";
</script>

<AuthRoute roles={[ UserRole.USER, UserRole.ADMIN ]}>
    <div class="flex flex-row justify-center">
        <slot/>
    </div>
</AuthRoute>