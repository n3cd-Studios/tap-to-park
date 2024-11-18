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
    <form class="p-10 bg-white rounded-xl w-1/3" on:submit|preventDefault={handleRegister} aria-labelledby="register-form">
        <Input bind:value={email} name="Email" id="email" aria-label="Email address"/>
        <Input bind:value={password} name="Password" type="password" id="password" aria-label="Password"/>
        {#if showInvite}
            <div aria-live="polite">
                <Input bind:value={invite} name="Invite Code" id="invite-code" aria-label="Invite Code"/>
            </div>
        {/if}
        <div class="flex flex-row justify-between">
            <button
                type="button"
                class="text-blue-400"
                aria-expanded={showInvite}
                aria-controls="invite-code"
                on:click={() => showInvite = !showInvite}>
                {`I ${showInvite ? "don't":""} have an invite code`}
            </button>
            <Button type="submit" aria-label="Register">Register</Button>
        </div>
    </form>
</div>