<template>
  <div class="loading-dashboard" role="status" aria-live="polite">
    <span class="loading-dashboard-mark">
      <img src="/img/logo.svg" alt="" />
    </span>
    <v-progress-linear
      indeterminate
      rounded
      color="primary"
      height="4"
      class="loading-dashboard-bar"
    />
    <p class="text-body-medium text-medium-emphasis mt-4 mb-0">
      {{ $t('components.loadingDashboard.loadingText') }}
    </p>
  </div>
</template>

<style scoped>
.loading-dashboard {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 100dvh;
  padding: 1rem;
  background: radial-gradient(
    circle at 50% 40%,
    rgb(var(--v-theme-primary), 0.12),
    transparent 40%
  );
}

.loading-dashboard-mark {
  display: inline-flex;
  width: 72px;
  height: 72px;
  margin-bottom: 1.5rem;
  padding: 0.6rem;
  border: 1px solid rgb(var(--v-border-color), var(--v-border-opacity));
  border-radius: 1.25rem;
  background: rgb(var(--v-theme-surface));
  box-shadow: 0 20px 50px rgb(var(--v-theme-primary), 0.25);
}

.loading-dashboard-mark img {
  width: 100%;
  height: 100%;
}

.loading-dashboard-bar {
  max-width: 180px;
}
</style>
