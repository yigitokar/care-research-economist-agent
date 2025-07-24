# Large Sample Theory: Analytical Foundations for Econometric Inference

> "The fundamental problem of statistical inference is that we see the sample but care about the population. Large sample theory is the mathematical bridge between these two worlds." — Whitney Newey (paraphrased)

## 1. The Econometrician's Fundamental Problem

Let $(\Omega, \mathcal{F}, P)$ be a probability space, and consider a sequence of i.i.d. random vectors $\{X_i\}_{i=1}^n$ drawn from distribution $P$. Our parameter of interest is $\theta_0 \in \Theta \subseteq \mathbb{R}^k$, defined implicitly as the unique solution to:

$$E_P[\psi(X_i, \theta)] = 0$$

where $\psi: \mathcal{X} \times \Theta \to \mathbb{R}^k$ is our moment function. This framework nests:
- **Mean estimation**: $\psi(X_i, \theta) = X_i - \theta$ with $\theta = E[X_i]$
- **Linear regression**: $\psi(W_i, Y_i, \theta) = (Y_i - W_i'\theta)W_i$
- **MLE**: $\psi(X_i, \theta) = \nabla_\theta \log f(X_i; \theta)$

**Key Insight**: The population parameter $\theta_0$ is defined by the infinite population, but we only observe the finite sample. Large sample theory studies how $\hat{\theta}_n \to \theta_0$ as $n \to \infty$.

## 2. Modes of Convergence: A Hierarchy

### 2.1 Convergence in Probability: The Weak Law

**Definition**: $\hat{\theta}_n \xrightarrow{p} \theta$ if for every $\epsilon > 0$:
$$P(\|\hat{\theta}_n - \theta\| > \epsilon) \to 0 \text{ as } n \to \infty$$

**The Critical Insight**: This is not saying $\hat{\theta}_n$ gets "close" to $\theta$ for large $n$. Rather, the *probability* of being far from $\theta$ becomes small. This is a statement about the *distribution* of $\hat{\theta}_n$, not $\hat{\theta}_n$ itself.

**Example (Sample Mean)**: Let $\bar{X}_n = \frac{1}{n}\sum_{i=1}^n X_i$ with $E[X_i] = \mu$ and $\text{Var}(X_i) = \sigma^2 < \infty$. By Chebyshev's inequality:

$$P(|\bar{X}_n - \mu| > \epsilon) \leq \frac{\sigma^2}{n\epsilon^2} \to 0$$

**Economic Intuition**: Consider estimating the average treatment effect (ATE) of a job training program. Even with large samples, extreme outcomes can occur, but their probability vanishes. This is why we can trust RCT results with sufficiently large samples.

### 2.2 Almost Sure Convergence: The Strong Law

**Definition**: $\hat{\theta}_n \xrightarrow{a.s.} \theta$ if:
$$P(\{\omega: \hat{\theta}_n(\omega) \to \theta\}) = 1$$

**The Deeper Insight**: Almost sure convergence is about the *trajectory* of the estimator across sequences of data. While convergence in probability allows for occasional large deviations (even as their probability shrinks), almost sure convergence requires that these deviations eventually cease.

**Kolmogorov's Strong Law**: If $\{X_i\}$ are i.i.d. with $E[|X_i|] < \infty$, then $\bar{X}_n \xrightarrow{a.s.} E[X_i]$.

**Application**: In time-series econometrics, stationarity conditions ensure that sample moments converge almost surely to population moments, justifying the use of sample autocorrelations in ARMA modeling.

### 2.3 Convergence in Distribution: The Central Limit Theorem

**Definition**: $\hat{\theta}_n \xrightarrow{d} \theta$ if:
$$P(\hat{\theta}_n \leq t) \to P(\theta \leq t)$$
for all continuity points $t$ of the CDF of $\theta$.

**Lindeberg-Lévy CLT**: Let $\{X_i\}$ be i.i.d. with $E[X_i] = \mu$ and $\text{Var}(X_i) = \sigma^2 < \infty$. Then:
$$\sqrt{n}(\bar{X}_n - \mu) \xrightarrow{d} N(0, \sigma^2)$$

**The Profound Implication**: The shape of the distribution of $\bar{X}_n$ becomes normal *regardless* of the underlying distribution of $X_i$. This is not a statement about the data—it is a statement about the *estimator's sampling distribution*.

**Economic Application**: When estimating demand elasticities with microdata, the distribution of individual preferences (log-normal, gamma, etc.) becomes irrelevant for inference on the *mean* elasticity as the sample grows.

## 3. The Delta Method: Propagating Uncertainty

### 3.1 Univariate Case

Let $\hat{\theta}_n \xrightarrow{p} \theta_0$ and $\sqrt{n}(\hat{\theta}_n - \theta_0) \xrightarrow{d} N(0, \Sigma)$. For a continuously differentiable function $g: \mathbb{R}^k \to \mathbb{R}$,:

$$\sqrt{n}(g(\hat{\theta}_n) - g(\theta_0)) \xrightarrow{d} N(0, \nabla g(\theta_0)' \Sigma \nabla g(\theta_0))$$

**Example**: Consider the elasticity of substitution $\sigma = \frac{\partial \log(C_1/C_2)}{\partial \log(p_1/p_2)}$ in a CES demand system. If we estimate $\hat{\theta}_n$ (the CES parameters), the delta method gives us the asymptotic distribution of $\hat{\sigma}_n = \sigma(\hat{\theta}_n)$.

### 3.2 Multivariate Delta Method

For $g: \mathbb{R}^k \to \mathbb{R}^m$:
$$\sqrt{n}(g(\hat{\theta}_n) - g(\theta_0)) \xrightarrow{d} N(0, G(\theta_0)' \Sigma G(\theta_0))$$
where $G(\theta_0) = \frac{\partial g}{\partial \theta'}(\theta_0)$.

**Application**: In structural estimation, we often estimate parameters $\theta$ and then compute counterfactual welfare measures $W(\theta)$. The delta method provides valid standard errors for these welfare calculations.

## 4. M-Estimation: A General Framework

### 4.1 Consistency via Uniform Law of Large Numbers

Let $\hat{\theta}_n = \arg\max_{\theta \in \Theta} \frac{1}{n}\sum_{i=1}^n m(X_i, \theta)$. Define:
- $Q_n(\theta) = \frac{1}{n}\sum_{i=1}^n m(X_i, \theta)$ (sample objective)
- $Q_0(\theta) = E[m(X_i, \theta)]$ (population objective)

**Theorem (Consistency)**: Under regularity conditions (compact $\Theta$, continuity, identification, and uniform convergence):
$$\hat{\theta}_n \xrightarrow{p} \theta_0 \equiv \arg\max_{\theta \in \Theta} Q_0(\theta)$$

**The Key Insight**: We need uniform convergence $\sup_{\theta \in \Theta} |Q_n(\theta) - Q_0(\theta)| \xrightarrow{p} 0$. This is much stronger than pointwise convergence and requires empirical process theory.

### 4.2 Asymptotic Normality

**Theorem**: Under additional regularity conditions (twice continuous differentiability, non-singular Hessian):
$$\sqrt{n}(\hat{\theta}_n - \theta_0) \xrightarrow{d} N(0, H^{-1} J H^{-1})$$
where:
- $H = E[\nabla_{\theta\theta'} m(X_i, \theta_0)]$ (expected Hessian)
- $J = E[\nabla_\theta m(X_i, \theta_0) \nabla_\theta m(X_i, \theta_0)']$ (outer product of gradients)

**Special Case - MLE**: When $m(X_i, \theta) = \log f(X_i; \theta)$, we have the information matrix equality: $H = -J$, so:
$$\sqrt{n}(\hat{\theta}_n - \theta_0) \xrightarrow{d} N(0, \mathcal{I}(\theta_0)^{-1})$$

**Economic Example**: Consider estimating a demand system with $J$ products. The M-estimator framework gives us the asymptotic distribution of price elasticities $\epsilon_{jk}(\hat{\theta}_n)$ for all products $j,k$.

## 5. GMM: When Moments Define Parameters

### 5.1 The GMM Framework

Let $\hat{\theta}_n = \arg\min_{\theta \in \Theta} \left[\frac{1}{n}\sum_{i=1}^n g(X_i, \theta)\right]' W_n \left[\frac{1}{n}\sum_{i=1}^n g(X_i, \theta)\right]$ where:
- $g: \mathcal{X} \times \Theta \to \mathbb{R}^q$ are moment conditions
- $W_n \xrightarrow{p} W$ is a positive definite weight matrix
- $E[g(X_i, \theta_0)] = 0$ (population moment conditions)

**Consistency**: $\hat{\theta}_n \xrightarrow{p} \theta_0$ under regularity conditions.

### 5.2 Asymptotic Distribution

$$\sqrt{n}(\hat{\theta}_n - \theta_0) \xrightarrow{d} N(0, (G'WG)^{-1}G'W\Omega WG(G'WG)^{-1})$$
where:
- $G = E[\nabla_\theta g(X_i, \theta_0)']$
- $\Omega = E[g(X_i, \theta_0)g(X_i, \theta_0)']$

**Efficient GMM**: When $W = \Omega^{-1}$, the asymptotic variance simplifies to:
$$\sqrt{n}(\hat{\theta}_n - \theta_0) \xrightarrow{d} N(0, (G'\Omega^{-1}G)^{-1})$$

**Application**: Estimating Euler equations in consumption-based asset pricing models. The moments are:
$$E[\beta R_{t+1}^{-\gamma}(C_{t+1}/C_t)^{-\gamma} | I_t] = 1$$
where $\theta = (\beta, \gamma)$ are the preference parameters.

## 6. Bootstrap and Resampling Methods

### 6.1 The Bootstrap Principle

Given data $\{X_i\}_{i=1}^n$, the nonparametric bootstrap samples with replacement to create $\{X_i^*\}_{i=1}^n$ and computes $\hat{\theta}_n^*$. The distribution of $\hat{\theta}_n^*$ approximates the sampling distribution of $\hat{\theta}_n$.

**Key Result**: Under regularity conditions, the bootstrap distribution is *asymptotically valid*:
$$\sup_t |P^*(\sqrt{n}(\hat{\theta}_n^* - \hat{\theta}_n) \leq t) - P(\sqrt{n}(\hat{\theta}_n - \theta_0) \leq t)| \xrightarrow{p} 0$$

**Economic Application**: In structural IO models, the bootstrap provides confidence intervals for counterfactual prices after mergers, where analytical standard errors are intractable.

### 6.2 Bootstrap Validity Conditions

The bootstrap requires:
1. **Smoothness**: The estimator must be sufficiently smooth (Hadamard differentiable)
2. **Moment conditions**: Sufficient moments must exist
3. **No boundary issues**: The true parameter must be in the interior of the parameter space

**Example**: The bootstrap fails for the maximum likelihood estimator of a parameter on the boundary (e.g., a variance component equal to zero).

## 7. High-Dimensional and Non-Regular Cases

### 7.1 When p → ∞ with n

In modern applications, the number of parameters $p$ grows with $n$. Consider:
- **Many instruments**: $p/n \to c \in (0,1)$
- **High-dimensional estimation**: $p \gg n$ with sparsity assumptions

**Key Results**:
- **Bekker (1994)**: Many weak instruments bias 2SLS toward OLS
- **LASSO**: Under sparsity, $\|\hat{\theta}_n - \theta_0\|_2 = O_P(\sqrt{s \log p/n})$ where $s$ is the sparsity level

### 7.2 Non-Regular Cases

**Examples**:
- **Parameter on boundary**: $\theta_0 \in \partial \Theta$
- **Non-differentiable moments**: Quantile regression with $m(X_i, \theta) = (Y_i - X_i'\theta)(\tau - 1\{Y_i \leq X_i'\theta\})$
- **Discontinuities**: Regression discontinuity designs

**Theory**: These cases require specialized limit theory (Chernoff, cube root asymptotics, etc.).

## 8. Practical Considerations and Common Pitfalls

### 8.1 The Finite-Sample Problem

Large sample theory provides *approximations*. The key question is: when are these approximations good enough?

**Edgeworth Expansions**: Provide higher-order approximations:
$$P(\sqrt{n}(\hat{\theta}_n - \theta_0) \leq t) = \Phi(t/\sigma) + \frac{\text{skewness}}{6\sqrt{n}}\Phi'''(t/\sigma) + O(1/n)$$

**Bootstrap vs. Analytical Standard Errors**: The bootstrap automatically captures higher-order terms but can be anti-conservative in small samples.

### 8.2 The Clustering Problem

With clustered data $\{(Y_{ig}, X_{ig})\}_{g=1}^G$, the effective sample size is $G$, not $n$. The asymptotic theory requires:
- **Many clusters**: $G \to \infty$
- **Bounded cluster sizes**: $\max_g n_g = O(1)$

**Cameron, Gelbach, and Miller (2008)**: Cluster-robust standard errors can be severely downward biased with few clusters.

## 9. The Frontier: Recent Developments

### 9.1 Double Machine Learning

For partially linear models:
$$Y = D\theta_0 + g_0(X) + U, \quad E[U|D,X] = 0$$

The DML estimator achieves:
$$\sqrt{n}(\hat{\theta}_{DML} - \theta_0) \xrightarrow{d} N(0, V)$$

even when $g_0(\cdot)$ is estimated by ML methods, provided the ML estimators satisfy certain rate conditions.

### 9.2 Causal Inference with High-Dimensional Controls

Recent work combines:
- **Approximate sparsity**: $\|\beta_0\|_0 \leq s \ll n$
- **Cross-fitting**: Sample splitting to avoid overfitting bias
- **Debiased estimators**: Constructing estimators that are regular (root-n consistent)

## 10. Summary: The Role of Large Sample Theory in Economic Research

Large sample theory provides the mathematical foundation for understanding:
1. **When estimators work**: Consistency conditions
2. **How well they work**: Asymptotic distributions and efficiency bounds
3. **How to quantify uncertainty**: Valid standard errors and inference
4. **The limits of inference**: When finite samples may not suffice

The theory is not merely technical—it shapes how we think about the relationship between observed data and unobserved population parameters that drive economic behavior.

## Mathematical Appendix

### A.1 Useful Inequalities

**Hölder's Inequality**: For $p,q > 1$ with $1/p + 1/q = 1$:
$$E[|XY|] \leq (E[|X|^p])^{1/p}(E[|Y|^q])^{1/q}$$

**Jensen's Inequality**: For convex $\phi$:
$$\phi(E[X]) \leq E[\phi(X)]$$

### A.2 Empirical Process Theory

**Donsker's Theorem**: The empirical process converges to a Brownian bridge:
$$\sqrt{n}(\hat{F}_n - F) \xrightarrow{d} \mathbb{B}(F)$$

### A.3 Stochastic Equicontinuity

A sequence of processes $\{G_n\}$ is stochastically equicontinuous if for every $\epsilon > 0$:
$$\lim_{\delta \to 0} \limsup_{n \to \infty} P^*(\sup_{\rho(s,t) < \delta} |G_n(s) - G_n(t)| > \epsilon) = 0$$

This condition is crucial for establishing uniform convergence results and asymptotic distributions of M-estimators.