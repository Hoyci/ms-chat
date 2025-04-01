import { SignupPayload, SignupSchema } from "@api/auth/types";
import { zodResolver } from "@hookform/resolvers/zod";
import { useAuthStore } from "@store/authStore";
import { Link, useNavigate } from "react-router-dom";
import { useForm } from "react-hook-form";

export function Signup() {
  const navigate = useNavigate();
  const { signup, loading, error } = useAuthStore();
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<SignupPayload>({
    resolver: zodResolver(SignupSchema),
  });

  const onSubmit = async (data: SignupPayload) => {
    await signup(data);
    navigate("/");
  };

  return (
    <div className="h-screen w-full bg-background flex items-center justify-center">
      <div className="w-96 bg-primary-100 p-8 rounded-lg shadow-lg">
        <h2 className="text-2xl font-bold text-white mb-6">Cadastro</h2>

        {error && (
          <div className="mb-4 p-2 bg-red-500 text-white text-sm rounded">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div>
            <input
              type="text"
              placeholder="Nome"
              className="w-full p-2 rounded bg-primary-200 text-white placeholder-gray-400"
              {...register("username")}
            />
            {errors.username && (
              <p className="text-red-500 text-xs mt-1">
                {errors.username.message}
              </p>
            )}
          </div>

          <div>
            <input
              type="email"
              placeholder="Email"
              className="w-full p-2 rounded bg-primary-200 text-white placeholder-gray-400"
              {...register("email")}
            />
            {errors.email && (
              <p className="text-red-500 text-xs mt-1">
                {errors.email.message}
              </p>
            )}
          </div>

          <div>
            <input
              type="password"
              placeholder="Senha"
              className="w-full p-2 rounded bg-primary-200 text-white placeholder-gray-400"
              {...register("password")}
            />
            {errors.password && (
              <p className="text-red-500 text-xs mt-1">
                {errors.password.message}
              </p>
            )}
          </div>

          <div>
            <input
              type="password"
              placeholder="Confirme sua senha"
              className="w-full p-2 rounded bg-primary-200 text-white placeholder-gray-400"
              {...register("confirmPassword")}
            />
            {errors.confirmPassword && (
              <p className="text-red-500 text-xs mt-1">
                {errors.confirmPassword.message}
              </p>
            )}
          </div>

          <button
            type="submit"
            disabled={loading}
            className="w-full bg-secondary-100 text-white py-2 rounded hover:bg-secondary-200 transition disabled:opacity-50 hover:cursor-pointer"
          >
            {loading ? "Criando conta..." : "Criar conta"}
          </button>
        </form>

        <div className="mt-4 text-center text-neutral-300">
          Já tem conta?{" "}
          <Link to="/login" className="text-secondary-100 hover:underline">
            Faça login
          </Link>
        </div>
      </div>
    </div>
  );
}
