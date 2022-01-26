<script setup lang="ts">
import { computed } from 'vue';
import { user } from '../stores/user'

const status = computed(() => {
    if (!user.fecthed) {
        return "Loading..."
    }

    if (user.isLoggedIn) {
        return "Logged in as " + user.email
    }

    return "Not logged in"
})
</script>

<template>
    <p class="auth-status">{{ status }}</p>
</template>

<style scoped>
/* place at the bottom, gray text on gray background */
.auth-status {
    background: #f8f8f8;
    border-radius: 3px;
    color: #999;
    font-size: 9px;
    padding: 3px 5px;
}
</style>
