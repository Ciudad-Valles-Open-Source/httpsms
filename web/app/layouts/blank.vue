<script setup lang="ts">
useThemeScope('public')
</script>

<template>
  <VApp>
    <VMain>
      <AppToast />
      <slot />
    </VMain>
  </VApp>
</template>
