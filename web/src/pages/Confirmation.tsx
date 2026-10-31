import { getRouteApi, useNavigate } from "@tanstack/react-router"
import { API_URL } from "../App"

const routeApi = getRouteApi("/confirm/$token")

export const Confirmation = () => {
    const { token } = routeApi.useParams()
    const navigate = useNavigate()

    const handleConfirm = async () => {
        const response = await fetch(`${API_URL}/users/activate`, {
            method: "PUT",
            body: JSON.stringify({
                token: token
            })
        })

        if (response.ok) {
            navigate({ to: "/" })
        } else {
            alert("Failed to confirm registeration")
        }
    }

    return (
        <div>
            <h1>Confirmation</h1>
            <button onClick={handleConfirm}>Click to confirm</button>
        </div>
    )
}