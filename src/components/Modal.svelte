<!-- modified https://tailwindui.com/components/application-ui/overlays/modal-dialogs -->
<script lang="ts">
    import Fa from "svelte-fa";
    import { fly } from "svelte/transition";
    import { faCheck, faRemove, faTimes } from "@fortawesome/free-solid-svg-icons";
    import { createEventDispatcher } from 'svelte';

    export let visible = false;
    export let title: string | null = null;
    export let message: string | null = null;
    export let error = false;
    export let buttonText: string | null = null;
    export let buttonFunctionality: (() => void) | null = null;

    const dispatch = createEventDispatcher();
    const closeModal = () => {
        dispatch('close');
    }
    
</script>

{#if visible}
    <div class="relative" aria-labelledby="modal-title" role="dialog" aria-modal="true">
        <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" aria-hidden="true"></div>
    
        <div class="fixed inset-0 w-screen overflow-y-auto">
            <div class="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
                <div class="relative transform overflow-hidden rounded-lg bg-white text-left shadow-xl transition-all sm:my-8 sm:w-full sm:max-w-lg"
                    in:fly={{ y: 40 }} out:fly={{ y: -40 }}>
                    <button 
                        class="absolute top-2 right-2 rounded-full text-gray-400 hover:text-gray-600 focus:outline-none"
                        on:click={closeModal}
                    >
                        <Fa icon={faTimes} class="h-5 w-5"/>
                    </button>
                    <div class="bg-white px-4 pb-4 pt-5 sm:p-6 sm:pb-4">
                        <div class="sm:flex sm:items-start">
                            <div class={`${!error ? 'bg-green-100' : 'bg-red-100'} mx-auto flex h-12 w-12 flex-shrink-0 items-center justify-center rounded-full sm:mx-0 sm:h-10 sm:w-10`}>
                                <slot name="icon">
                                    {#if !error}
                                        <Fa icon={faCheck} class="text-green-600 text-2xl"/>
                                    {:else}
                                        <Fa icon={faRemove} class="text-red-600 text-2xl"/>  
                                    {/if}
                                </slot>
                            </div>
                        <div class="mt-3 text-center sm:ml-4 sm:mt-0 sm:text-left">
                            <slot name="title">
                                {#if title}
                                    <h3 class="text-base font-semibold leading-6 text-gray-900" id="modal-title">{title}</h3>
                                {/if}
                            </slot>
                            <slot name="message">
                                {#if message}
                                    <div class="mt-2">
                                        <p class="text-sm text-gray-500">{message}</p>
                                    </div>
                                {/if}
                            </slot>
                        </div>
                    </div>
                </div>
                <slot name="button">
                    {#if buttonText && buttonFunctionality}
                        <div class="bg-gray-50 px-4 py-3 sm:flex sm:flex-row-reverse sm:px-6">
                            <button on:click={buttonFunctionality} type="button" class={`${!error ? 'bg-green-600 hover:bg-green-500' : 'bg-red-600 hover:bg-red-500'} inline-flex w-full justify-center rounded-md px-3 py-2 text-sm font-semibold text-white shadow-sm sm:ml-3 sm:w-auto`}>{buttonText}</button>
                        </div>
                    {/if}
                </slot>
                </div>
            </div>
        </div>
    </div>
{/if}