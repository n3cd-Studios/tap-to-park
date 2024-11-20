<script>
    import { UserRole } from "$lib/models";
    import Hover from "../../components/animation/Hover.svelte";
import AuthRoute from "../../components/auth/AuthRoute.svelte";
</script>

<AuthRoute roles={[ UserRole.ADMIN ]} redirect="/user">
    <div class="flex flex-col justify-around mt-10 sm:flex-row gap-2" role="main" aria-label="Admin dashboard">
        <div class="flex flex-col justify-between mx-3 sm:ml-10 bg-white pl-6 pr-14 py-3 rounded-t-2xl" aria-label="Admin navigation">
            <div class="flex flex-col gap-2" role="menu" aria-label="Admin menu options">
                <Hover href="/admin" aria-label="Navigate to admin dashboard">Dashboard</Hover>
                <Hover href="/admin/invites" aria-label="Navigate to invite management">Invites</Hover>
                <Hover href="/admin/analytics" aria-label="Navigate to analytics dashboard">Analytics</Hover>
            </div>
            <div role="menuitem" aria-label="Logout">
                <Hover href="/auth/logout" aria-label="Log out of your account">Logout</Hover>
            </div>
        </div>
        <div class="mx-3 sm:w-full sm:mx-10" role="region" aria-label="Admin page content">
            <slot />
        </div>
    </div>
</AuthRoute>
