<template>
  <v-app :theme="layoutTheme">
    <Onboarding @wrapper="reference => wrapper = reference" :steps="steps" @finish="onboarding = false" :close="finish" />
    <AppBar>
      <template #prepend-actions>
        <v-app-bar-nav-icon @click="drawer = !drawer" id="example-menu"></v-app-bar-nav-icon>
      </template>
      <template #desktop-actions>
        <PrimeButton @click="start" v-if="onboarding" variant="outlined">Onboarding</PrimeButton>
        <ShareButton btn-class="ml-2" />
        <SaveButton btn-class="ml-2" />
        <LoadButton btn-class="mx-2" />
      </template>

      <template #mobile-actions>
        <MobileMenu v-model:policy="inputs.policy" v-model:resource="inputs.resource" v-model:context="inputs.context" />
      </template>
    </AppBar>
    <ExampleDrawer v-model="drawer" />
    <v-main>
      <v-container fluid class="pr-lg-4 pr-md-4 pr-sm-8 pr-8">
        <OnboardingAlert />
        <v-row>
          <v-col :md="7" :sm="12">
            <PolicyPanel v-model="inputs.policy" />
          </v-col>
          <v-col :md="5" :sm="12">
            <ContextPanel v-model="inputs.context" @collapse="(v: boolean) => resourceHeight = v ? 691 : 441" />
            <ResourcePanel v-model="inputs.resource" :height="resourceHeight" />
          </v-col>
        </v-row>
        <HelpButton />
        <StartButton
          @on-response="handleResponse"
          @on-error="handleError"
          :error-state="showError"
        />
      </v-container>
      <ErrorBar v-model="showError" :text="errorText" />
      <ResultDialog v-model="showResults" :results="results" :policy="inputs.policy" />
      <v-card v-if="state.name.value" class="state-card">
        <v-card-text class="text-body-2 font-weight-medium py-2">Loaded State: {{ state.name.value }}</v-card-text>
      </v-card>
    </v-main>
    <Sponsor />
    <AppVersion />
  </v-app>
</template>

<script setup lang="ts">
import { ref, watchEffect, onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import 'v-onboarding/dist/style.css'

import { layoutTheme } from "@/config";
import { useState, useOnboarding } from "@/composables";
import { inputs, populate } from "@/store"
import { EngineResponse } from "@/types";
import ErrorBar from "@/components/ErrorBar.vue";
import ExampleDrawer from "@/components/ExampleDrawer.vue";
import OnboardingAlert from "@/components/OnboardingAlert.vue";
import HelpButton from '@/components/HelpButton.vue';
import Onboarding from "@/components/Onboarding.vue";
import Sponsor from "@/components/Sponsor.vue";
import PrimeButton from "@/components/PrimeButton.vue";
import { ResourcePanel, ContextPanel, PolicyPanel } from "@/components/Panel";
import { LoadButton, SaveButton, ShareButton, MobileMenu, AppBar } from "@/components/AppBar";
import { StartButton, ResultDialog } from "@/components/Results";
import { parseContent } from "@/functions/share";
import AppVersion from "@/components/AppVersion.vue";
import { loadFromRepo } from "@/functions/github";

const route = useRoute();
const router = useRouter();

const state = useState()

const showResults = ref<boolean>(false);
const results = ref<EngineResponse>({ policies: [], resources: [] });

const handleResponse = (response: EngineResponse) => {
  results.value = response;
  showResults.value = true;
};

const showError = ref<boolean>(false);
const errorText = ref<string>("");

const handleError = (error: Error) => {
  errorText.value = error.message;
  showError.value = true;

  setTimeout(() => {
    errorText.value = "";
    showError.value = false;
  }, 5000);
};

const drawer = ref<boolean>(false);

watchEffect(() => {
  const query = route.query.content as string;
  if (query) {
    try {
      parseContent(query)

      router.replace({ ...route, query: {} });
      return
    } catch (err) {
      console.error("could not parse content string", err);
    }
  }

  const policy = route.query.policy as string;
  if (policy) {
    loadFromRepo(policy, route.query.resource as string).finally(() => {
      router.replace({ ...route, query: {} });
    })
  }
});

onMounted(() => {
  if (route.query.content) {
    return;
  }

  populate()
});

const { finish, start, onboarding, steps, wrapper } = useOnboarding(drawer)

const resourceHeight = ref(441)
</script>

<style scoped>
.state-card {
  background-color: rgb(var(--v-theme-background));
  position: fixed; 
  bottom: 0; 
  left: 0;
}
</style>