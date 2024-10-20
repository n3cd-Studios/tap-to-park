<script lang="ts">
    import { goto } from "$app/navigation";
    import { register } from "$lib/auth";
    import { UserRole } from "$lib/models";
    import Button from "../../../components/form/Button.svelte";
    import Input from "../../../components/form/Input.svelte";
    import { toaster } from "../../../components/toaster/toaster";

    let email: string;
    let password: string;
    let invite: string;

    let showInvite = false;

    const handleRegister = async () => {
    await register(email, password, invite)
        .then(user => user?.role == UserRole.ADMIN ? goto("/admin") : goto("/user"))
        .catch(() => toaster.push({ type: "error", message: "Failed to register." }, 5000));
    };

</script>

<div class="h-full w-full flex flex-col justify-center items-center">
    <form class="p-10 bg-white rounded-xl w-1/3" on:submit|preventDefault={handleRegister}>
        <Input bind:value={email} name="Email"/>
        <Input bind:value={password} name="Password" type="password"/>
        {#if showInvite}
            <Input bind:value={invite} name="Invite Code" />
        {/if}
        <div class="flex flex-row justify-between">
            <button type="button" class="text-blue-400" on:click={() => showInvite = !showInvite}>
                {`I ${showInvite ? "don't":""} have an invite code`}
            </button>
            <Button type="submit">Register</Button>
        </div>
    </form>
</div>