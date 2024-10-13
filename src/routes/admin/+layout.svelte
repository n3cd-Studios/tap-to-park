<script>
    import Hover from "../../components/animation/Hover.svelte";
import AuthRoute from "../../components/auth/AuthRoute.svelte";
</script>

<AuthRoute>
    <div class="flex flex-row justify-around mt-10">
        <div class="flex flex-col gap-2 ml-10 bg-white pl-6 pr-14 py-3 rounded-lg">
            <Hover href="/admin">Dashboard</Hover>
            <Hover href="/admin/invites">Invites</Hover>
        </div>
        <div class="w-full mx-10">
            <slot />
        </div>
    </div>
</AuthRoute>