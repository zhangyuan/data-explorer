<script setup lang="ts"></script>

<template>
    <div class=" bg-gray-200">
        <nav class="container mx-auto p-3 text-center text-2xl">
            <ul class="flex justify-center text-black">
                <li>
                    <RouterLink to="/" class="hover:underline underline-offset-4">主页</RouterLink>
                </li>
                <li class="mx-3">
                    <RouterLink to="/issues" class="hover:underline underline-offset-4">议题</RouterLink>
                </li>
            </ul>
        </nav>
        <main class="container mx-auto">
            <RouterView />
        </main>
    </div>
</template>

<style scoped></style>
