import { Then } from 'cucumber';
const tasksSteps = require(__srcdir + '/steps/tasks/tasksSteps.js');

let tSteps = new tasksSteps(__wdriver);

Then(/^the Tasks page is loaded$/, {timeout: 2 * 5000}, async() => {
    await tSteps.isLoaded();
    await tSteps.verifyIsLoaded();
    await tSteps.verifyHeaderContains('Tasks', 'tasks-page');
});

Then(/^there is a task named "(.*)"$/, async name =>{
    await tSteps.verifyTaskCardVisible(name);
});
