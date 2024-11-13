<script lang="ts">
    import Fa from 'svelte-fa';
    import { faBars } from '@fortawesome/free-solid-svg-icons';
    import Button from '../form/Button.svelte';
    import { ButtonType } from "$lib/utils";
    import { getUserInfo } from "$lib/auth";
    import { onMount } from "svelte";
    import { UserRole } from '$lib/models';
  
    export let onLoginRedirect = '/auth/login';
  
    let isLoggedIn = false;
    let userEmail: string | null = null;
    let isAdmin = false;
    let dropdownOpen = false;
    let container: HTMLDivElement;
  
    const generateDropdownOptions = (): { label: string, route: string }[] => {
        const options = [
            { label: 'User Profile', route: '/user' },
            { label: 'Settings & Privacy', route: '/auth/settings' },
            ...(isAdmin ? [{ label: 'Admin Dashboard', route: '/admin' }] : []),
            { label: 'Logout', route: '/auth/logout' }
        ];
  
        return options;
    };
  
    let dropdownOptions = generateDropdownOptions();
  
    onMount(async () => {
        try {
            const user = await getUserInfo();
            if (user && user.email) {
                userEmail = user.email;
                isLoggedIn = true;
                isAdmin = user.role === UserRole.ADMIN;
                dropdownOptions = generateDropdownOptions();
            } else {
                isLoggedIn = false;
            }
        } catch (error) {
            console.error("Failed to fetch user info:", error);
            isLoggedIn = false;
        }
    });
  
    const handleLoginButton = (event: MouseEvent) => {
        event.stopPropagation();
        if (isLoggedIn) {
            dropdownOpen = !dropdownOpen;
        } else {
            location.href = onLoginRedirect;
        }
    };
  
    const handleDropdownSelection = (route: string) => {
        location.href = route;
    };
  
    function onWindowClick(e: MouseEvent) {
        if (container && !container.contains(e.target as Node)) {
            dropdownOpen = false;
        }
    }

</script>
  
<svelte:window on:click={onWindowClick} />
<div bind:this={container}>
    <Button buttonType={ButtonType.DEFAULT} on:click={handleLoginButton}>
        {#if isLoggedIn}
            <span style="display: flex; align-items: center;">
                <Fa icon={faBars} class="mr-2 fa-bounce" /> {userEmail}
            </span>
        {:else}
            Login
        {/if}
    </Button>

    {#if dropdownOpen && isLoggedIn}
        <div class="absolute right-0 mt-2 w-48 rounded-md shadow-lg bg-white ring-1 ring-black ring-opacity-5">
            <div class="py-1" role="menu" aria-orientation="vertical" aria-labelledby="options-menu">
                {#each dropdownOptions as { label, route }}
                    <button on:click={() => handleDropdownSelection(route)} class="block px-4 py-2 text-med text-[#021427] hover:bg-gray-100 w-full text-left">
                        {label}
                    </button>
                {/each}
            </div>
        </div>
    {/if}
</div>

<style>
    .btn {
        /* Style your button accordingly */
    }
</style>
  