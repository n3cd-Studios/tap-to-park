<script>
    import { UserRole } from "$lib/models";
    import Hover from "../../components/animation/Hover.svelte";
import AuthRoute from "../../components/auth/AuthRoute.svelte";
</script>

<AuthRoute roles={[ UserRole.ADMIN ]} redirect="/user">
    <div class="flex flex-col justify-around mt-10 sm:flex-row gap-2">
        <div class="flex flex-col gap-2 mx-3 sm:ml-10 bg-white pl-6 pr-14 py-3 rounded-2xl">
            <Hover href="/admin">Dashboard</Hover>
            <Hover href="/admin/invites">Invites</Hover>
        </div>
        <div class="mx-3 sm:w-full sm:mx-10">
            <slot />
        </div>
    </div>
</AuthRoute>