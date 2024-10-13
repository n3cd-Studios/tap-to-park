<script>
    import Hover from "../../components/animation/Hover.svelte";
import AuthRoute from "../../components/auth/AuthRoute.svelte";
</script>

<AuthRoute>
    <div class="flex flex-row justify-around mt-10">
        <div class="flex flex-col mx-20 gap-2">
            <Hover href="/admin">Dashboard</Hover>
            <Hover href="/admin/invites">Invites</Hover>
        </div>
        <div class="w-1 bg-white"></div>
        <div class="w-full mx-10">
            <slot />
        </div>
    </div>
</AuthRoute>