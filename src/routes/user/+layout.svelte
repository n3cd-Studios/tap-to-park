<script lang="ts">
    import { type User, UserRole } from "$lib/models";
    import Hover from "../../components/animation/Hover.svelte";
    import AuthRoute from "../../components/auth/AuthRoute.svelte";
</script>

<Metadata title="User" description="Tap-To-Park user dashboard" />
<AuthRoute roles={[ UserRole.USER, UserRole.ADMIN ]}>
    <div class="flex flex-col justify-around mt-10 sm:flex-row gap-2" role="main" aria-label="User dashboard">
        <div class="flex flex-col justify-between mx-3 sm:ml-10 bg-white pl-6 pr-14 py-3 rounded-t-2xl" aria-label="User navigation">
            <div class="flex flex-col gap-2" role="menu" aria-label="User menu options">
                <Hover href="/user" aria-label="Navigate to user dashboard">Transactions</Hover>
                <Hover href="/user/sessions" aria-label="Navigate to session dashboard">Sessions</Hover>
            </div>
            <div role="menuitem" aria-label="Logout">
                <Hover href="/auth/logout" aria-label="Log out of your account">Logout</Hover>
            </div>
        </div>
        <div class="mx-3 sm:w-full sm:mx-10" role="region" aria-label="User page content">
            <slot />
        </div>
    </div>
</AuthRoute>
