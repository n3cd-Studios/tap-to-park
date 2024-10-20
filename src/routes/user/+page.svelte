
<script lang="ts">
    import type { User } from "$lib/models";

    export let data: User;
</script>

<p>Hello, <span class="font-bold text-center">{data.email}</span>. <a href="/auth/logout">Logout.</a></p>