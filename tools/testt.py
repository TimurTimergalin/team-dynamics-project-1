import os
import monorepo_root


print(*os.walk(os.path.join(monorepo_root.get_root(), "deployment")), sep='\n')
