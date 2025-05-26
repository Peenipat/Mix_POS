import type { NavigateFunction } from "react-router-dom";

export function navigateByRole(role: string | undefined, navigate: NavigateFunction) {
    if (!role) {
        navigate("/dashboard");
        return;
    }
    switch (role) {
        case 'SAAS_SUPER_ADMIN':
            navigate('/admin/dashboard');
            break;
        case 'TENANT':
            navigate('/tenant/dashboard');
            break;
        case 'BRANCH_ADMIN':
            navigate('/admin/dashboard');
            break;
        case 'STAFF':
            navigate('/staff/dashboard');
            break;
        default:
            navigate('/dashboard');
    }
}
