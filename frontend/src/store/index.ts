import { ConfigTemplate, ContextTemplate, PolicyTemplate, ResourceTemplate } from "@/assets/templates";
import { reactive } from "vue";
import { useState, Inputs } from "@/composables";

export const inputs = reactive({
    policy: PolicyTemplate,
    resource: ResourceTemplate,
    context: ContextTemplate,
    config: ConfigTemplate,
})

export const reset = () => {
    inputs.policy = PolicyTemplate
    inputs.resource = ResourceTemplate
    inputs.context = ContextTemplate
    inputs.context = ConfigTemplate
}

export const setDefaults = () => {
    init({ 
        policy: PolicyTemplate,
        resource: ResourceTemplate,
        context: ContextTemplate,
        config: ConfigTemplate,
    })
}

export const init = (values: Inputs) => {
    const state = useState()

    if (typeof values.policy === 'string') {
        state.policy.value = values.policy;
        inputs.policy = values.policy;
    }

    if (typeof values.resource === 'string') {
        state.resource.value = values.resource;
        inputs.resource = values.resource;
    }

    if (typeof values.context === 'string') {
        state.context.value = values.context;
        inputs.context = values.context;
    }

    if (typeof values.config === 'string') {
        state.config.value = values.config;
        inputs.config = values.config;
    }
    
    state.name.value = values.name || "";
}

export const update = (values: Inputs) => {
    if (values.policy) {
        inputs.policy = values.policy;
    }

    if (values.resource) {
        inputs.resource = values.resource;
    }

    if (values.context) {
        inputs.context = values.context;
    }
}

export const populate = () => {
    const state = useState()

    if (state.policy.value) {
        inputs.policy = state.policy.value;
    }
    if (state.resource.value) {
        inputs.resource = state.resource.value;
    }
    if (state.context.value) {
        inputs.context = state.context.value;
    }
    if (state.config.value) {
        inputs.config = state.config.value;
    }
}