<script setup lang="ts">
import ListController from '@/components/ListController.vue'
import { ref, watchEffect } from 'vue'
import { useLists } from '../composables/useLists'
import AuthStatus from '../components/AuthStatus.vue'
import { user } from '../stores/user'

const { lists } = useLists()
const selectedListID = ref('$meta')

watchEffect(() => {
    if (!user.isLoggedIn) {
        selectedListID.value = '$meta'
    }
})
</script>

<template>
    <!-- Side navigation contains header and links to all items -->
    <div>
        <div class="side-nav">
            <h3 class="side-header">rw youtube</h3>
            <div class="side-nav-flex border-top">
                <a
                    v-for="list in lists"
                    :key="list.ID"
                    :class="{ 'side-nav-selected': list.ID == selectedListID }"
                    @click="selectedListID = list.ID"
                >
                    {{ list.Name }}
                </a>
            </div>
            <AuthStatus class="auth-status" />
        </div>

        <!-- Main content located aside of navigation -->
        <div class="main-content">
            <ListController :list-id="selectedListID" />
        </div>
    </div>
</template>

<style scoped>
.auth-status {
    margin-top: auto;
}

.side-nav {
    background: #f8f8f8;
    border-right: 1px solid #e7e7e7;
    height: 100vh;
    overflow: hidden;
    position: fixed;
    top: 0;
    width: 200px;
    z-index: 1;
    display: flex;
    flex-direction: column;
}

.side-nav > h3 {
    font-weight: 1200;
    margin-left: 5%;
}

.side-nav-flex {
    display: flex;
    flex-direction: column;
    padding: 0;
    margin: 0;
}

.side-nav-flex > a {
    text-decoration: none;
    width: 100%;
    padding: 10px;
}

.side-nav-flex > a:hover,
.side-nav-selected {
    --tw-bg-opacity: 1;
    background-color: rgba(209, 213, 219, var(--tw-bg-opacity));
}

.border-top {
    border-top: 1px solid #e7e7e7;
}

.main-content {
    margin-left: 200px;
    padding: 20px;
}
</style>
