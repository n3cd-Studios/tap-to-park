
<script>
    import { authStore } from "$lib/auth";

    $: email = $authStore.user?.email
</script>

<p>Organization management</p>