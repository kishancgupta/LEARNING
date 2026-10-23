const descriptions = {
  bubble: {
    title: 'Bubble sort moves the biggest values to the end.',
    copy: 'It checks neighboring pairs and swaps them when they are out of order. Each pass settles one more value into its final position.',
  },
  selection: {
    title: 'Selection sort repeatedly chooses the smallest value.',
    copy: 'It scans the unsorted section to find its minimum, then swaps that value into the next open position.',
  },
  insertion: {
    title: 'Insertion sort grows a sorted section from left to right.',
    copy: 'Each new value travels left through the sorted section until it reaches the right place—like sorting a hand of cards.',
  },
  merge: {
    title: 'Merge sort divides first, then combines in order.',
    copy: 'It splits the list into smaller halves, sorts each half, and merges them by repeatedly choosing the smaller next value.',
  },
  quick: {
    title: 'Quick sort partitions values around a pivot.',
    copy: 'Values smaller than the pivot go left and larger values go right. It repeats that strategy on each section.',
  },
  linear: {
    title: 'Linear search checks each value in sequence.',
    copy: 'Starting at the beginning, it compares each item with the target until it finds a match or reaches the end.',
  },
  binary: {
    title: 'Binary search cuts the search range in half.',
    copy: 'It compares the middle value with the target and discards half the remaining values each turn. The input must be sorted.',
  },
};

const samples = {
  random: [42, 18, 65, 29, 73, 11, 54, 36, 82, 23, 61, 7],
  sorted: [7, 11, 18, 23, 29, 36, 42, 54, 61, 65, 73, 82],
  nearly: [12, 18, 21, 29, 34, 33, 41, 48, 52, 60, 67, 75],
  reversed: [88, 74, 63, 51, 42, 35, 27, 19, 12, 8],
  duplicates: [5, 3, 8, 3, 1, 8, 5, 2, 3, 1, 7, 5],
};

const modeButtons = [...document.querySelectorAll('.segment')];
const algorithmSelect = document.querySelector('#algorithm');
const datasetSelect = document.querySelector('#dataset');
const valuesInput = document.querySelector('#values');
const targetInput = document.querySelector('#target');
const targetField = document.querySelector('#target-field');
const chart = document.querySelector('#chart');
const emptyChart = document.querySelector('#chart-empty');
const errorMessage = document.querySelector('#error');
const progress = document.querySelector('#progress');
const state = { mode: 'sorting', algorithms: null, trace: null, step: 0, timer: null, original: [] };

function updateAlgorithmOptions() {
  const algorithms = state.algorithms?.[state.mode] ?? [];
  algorithmSelect.replaceChildren(...algorithms.map(({ id, name }) => {
    const option = document.createElement('option');
    option.value = id;
    option.textContent = name;
    return option;
  }));
  targetField.hidden = state.mode !== 'searching';
  updateLesson();
}

function updateLesson() {
  const algorithm = algorithmSelect.value;
  const info = descriptions[algorithm];
  const complexity = state.algorithms?.[state.mode]?.find((item) => item.id === algorithm)?.complexity ?? '—';
  document.querySelector('#complexity').textContent = complexity;
  document.querySelector('#lesson-title').textContent = info?.title ?? 'Every algorithm has a strategy.';
  document.querySelector('#lesson-copy').textContent = info?.copy ?? 'Select an algorithm to see what makes it unique.';
  document.querySelector('#lesson-complexity').textContent = complexity;
}

function setValues(values) {
  valuesInput.value = values.join(', ');
  state.original = [...values];
  renderBars(values, [], -1, -1, false);
  resetTrace();
}

function parseValues() {
  const raw = valuesInput.value.trim();
  const values = raw.split(',').map((item) => item.trim());
  if (!raw || values.some((item) => !/^-?\d+$/.test(item))) {
    throw new Error('Enter whole numbers separated by commas.');
  }
  const numbers = values.map(Number);
  if (numbers.length > 30) throw new Error('Use 30 values or fewer to keep each step easy to follow.');
  if (numbers.some((item) => !Number.isSafeInteger(item))) {
    throw new Error('Each value must be a safe whole number.');
  }
  return numbers;
}

function renderBars(values, indices, pivotIndex, low, high, found) {
  chart.replaceChildren();
  emptyChart.hidden = values.length > 0;
  const min = Math.min(...values);
  const max = Math.max(...values);
  const range = max - min || 1;
  values.forEach((value, index) => {
    const column = document.createElement('div');
    column.className = 'bar-column';
    const label = document.createElement('span');
    label.className = 'bar-value';
    label.textContent = value;
    const bar = document.createElement('div');
    bar.className = 'bar';
    bar.style.height = `${18 + ((value - min) / range) * 74}%`;
    if (found && indices.includes(index)) bar.classList.add('found');
    else if (index === pivotIndex) bar.classList.add('pivot');
    else if (indices.includes(index)) bar.classList.add('compare');
    else if (state.mode === 'searching' && algorithmSelect.value === 'binary' && low >= 0 && (index < low || index > high)) bar.classList.add('inactive');
    else if (state.mode === 'sorting' && found) bar.classList.add('found');
    const indexLabel = document.createElement('span');
    indexLabel.className = 'bar-index';
    indexLabel.textContent = index;
    bar.append(indexLabel);
    column.append(label, bar);
    chart.append(column);
  });
}

function updateStep() {
  if (!state.trace) return;
  const step = state.trace.steps[state.step];
  const done = state.step === state.trace.steps.length - 1 && step.action === 'done';
  const found = step.action === 'found';
  const middle = state.mode === 'searching' && algorithmSelect.value === 'binary';
  renderBars(step.values, step.indices ?? [], step.action === 'pivot' || middle ? step.indices?.[0] ?? step.mid : -1,
    step.low ?? -1, step.high ?? -1, found || done);
  document.querySelector('#message').textContent = step.message;
  document.querySelector('#step-count').textContent = `Step ${state.step + 1} of ${state.trace.steps.length}`;
  progress.max = String(state.trace.steps.length - 1);
  progress.value = String(state.step);
  const visibleSteps = state.trace.steps.slice(0, state.step + 1);
  document.querySelector('#comparisons').textContent = visibleSteps.filter((item) => item.action === 'compare').length;
  document.querySelector('#swaps').textContent = visibleSteps.filter((item) => item.action === 'swap').length;
  document.querySelector('#result').textContent = state.mode === 'sorting'
    ? (done ? 'Sorted' : '—')
    : found ? `Index ${state.trace.foundIndex}` : (step.action === 'not-found' ? 'Not found' : '—');
  document.querySelector('#run-title').textContent = `${algorithmSelect.options[algorithmSelect.selectedIndex]?.text ?? ''} · ${state.mode === 'sorting' ? 'Sorting' : 'Searching'}`;
  const playButton = document.querySelector('#play');
  playButton.textContent = state.step >= state.trace.steps.length - 1 ? '↺' : state.timer ? 'Ⅱ' : '▶';
  playButton.setAttribute('aria-label', state.timer ? 'Pause' : state.step >= state.trace.steps.length - 1 ? 'Replay' : 'Play');
}

function resetTrace() {
  stopPlayback();
  state.trace = null;
  state.step = 0;
  progress.max = '0';
  progress.value = '0';
  document.querySelector('#step-count').textContent = 'Step 0 of 0';
  document.querySelector('#message').textContent = 'Your step-by-step explanation will appear here.';
  document.querySelector('#comparisons').textContent = '0';
  document.querySelector('#swaps').textContent = '0';
  document.querySelector('#result').textContent = '—';
  document.querySelector('#run-title').textContent = 'Ready when you are';
}

function stopPlayback() {
  if (state.timer) clearInterval(state.timer);
  state.timer = null;
  if (state.trace) updateStep();
}

async function runAlgorithm() {
  stopPlayback();
  errorMessage.hidden = true;
  try {
    const values = parseValues();
    state.original = [...values];
    const searching = state.mode === 'searching';
    const target = Number(targetInput.value);
    if (searching && (!targetInput.value.trim() || !Number.isSafeInteger(target))) {
      throw new Error('Enter a whole-number target.');
    }
    if (searching && algorithmSelect.value === 'binary' && values.some((value, i) => i > 0 && value < values[i - 1])) {
      throw new Error('Binary search needs ascending data. Choose the “Sorted values” sample or sort your values first.');
    }
    const response = await fetch(searching ? '/api/search' : '/api/sort', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ algorithm: algorithmSelect.value, values, target }),
    });
    const result = await response.json();
    if (!response.ok) throw new Error(result.error ?? 'The algorithm could not be run.');
    state.trace = result;
    state.step = 0;
    updateStep();
  } catch (error) {
    errorMessage.textContent = error.message;
    errorMessage.hidden = false;
  }
}

function moveStep(direction) {
  if (!state.trace) return;
  state.step = Math.max(0, Math.min(state.trace.steps.length - 1, state.step + direction));
  updateStep();
}

function togglePlayback() {
  if (!state.trace) {
    runAlgorithm();
    return;
  }
  if (state.timer) {
    stopPlayback();
    return;
  }
  if (state.step >= state.trace.steps.length - 1) state.step = 0;
  updateStep();
  state.timer = setInterval(() => {
    if (state.step >= state.trace.steps.length - 1) {
      stopPlayback();
      return;
    }
    state.step++;
    updateStep();
  }, Number(document.querySelector('#speed').value));
  updateStep();
}

modeButtons.forEach((button) => button.addEventListener('click', () => {
  state.mode = button.dataset.mode;
  modeButtons.forEach((item) => item.classList.toggle('active', item === button));
  updateAlgorithmOptions();
  resetTrace();
}));
algorithmSelect.addEventListener('change', () => {
  updateLesson();
  resetTrace();
});
datasetSelect.addEventListener('change', () => {
  if (datasetSelect.value !== 'custom') setValues(samples[datasetSelect.value]);
});
valuesInput.addEventListener('input', () => {
  datasetSelect.value = 'custom';
  resetTrace();
});
document.querySelector('#run').addEventListener('click', runAlgorithm);
document.querySelector('#shuffle').addEventListener('click', () => {
  datasetSelect.value = 'random';
  setValues(Array.from({ length: 12 }, () => Math.floor(Math.random() * 90) + 5));
});
document.querySelector('#previous').addEventListener('click', () => moveStep(-1));
document.querySelector('#next').addEventListener('click', () => moveStep(1));
document.querySelector('#play').addEventListener('click', togglePlayback);
document.querySelector('#reset').addEventListener('click', () => {
  if (state.trace) {
    stopPlayback();
    state.step = 0;
    updateStep();
  } else {
    setValues(state.original);
  }
});
progress.addEventListener('input', () => {
  if (state.trace) {
    state.step = Number(progress.value);
    updateStep();
  }
});
document.querySelector('#speed').addEventListener('change', () => {
  document.querySelector('#speed-label').textContent = `${document.querySelector('#speed').selectedOptions[0].text} speed`;
  if (state.timer) {
    stopPlayback();
    togglePlayback();
  }
});

async function initialize() {
  try {
    const response = await fetch('/api/algorithms');
    if (!response.ok) throw new Error('Could not load the algorithm list.');
    state.algorithms = await response.json();
    updateAlgorithmOptions();
    setValues(samples.random);
  } catch (error) {
    errorMessage.textContent = error.message;
    errorMessage.hidden = false;
  }
}

initialize();
